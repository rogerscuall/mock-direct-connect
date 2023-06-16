package main

import (
	"dx-mock/adapters/db"
	d "dx-mock/pkg/dx"
	"encoding/json"
	"log"
	"net/http"
)

// CreateBGPPeer creates a BGP Peer.
// It checks if the Virtual Interface exists in the database and is available.
// If it is it will add the BGP Peer to the Virtual Interface.
func CreateBGPPeer(w http.ResponseWriter, r *http.Request) {
	var req d.CreateBgpPeerRequest
	err := d.RequestToJson(r, &req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vifDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer vifDB.CloseDbConnection()

	var privateVIF d.PrivateVirtualInterface

	err = vifDB.GetVal(req.VirtualInterfaceID, &privateVIF)
	if err != nil {
		log.Println("Error in getting virtual interface ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	// Check if the Virtual Interface is available
	if privateVIF.VirtualInterfaceState != "available" {
		http.Error(w, "Virtual Interface is not available", http.StatusBadRequest)
		return
	}

	// check if the BGP Peer already exists
	for _, bgpPeer := range privateVIF.BGPPeers {
		if bgpPeer.BGPPeerID == req.NewBGPPeer.CustomerAddress {
			http.Error(w, "BGP Peer already exists", http.StatusBadRequest)
			return
		}
	}

	// Add the BGP Peer to the Virtual Interface
	privateVIF.BGPPeers = append(privateVIF.BGPPeers, d.BGPConfig{
		AddressFamily:      req.NewBGPPeer.AddressFamily,
		AmazonAddress:      req.NewBGPPeer.AmazonAddress,
		ASN:                req.NewBGPPeer.Asn,
		AuthKey:            req.NewBGPPeer.AuthKey,
		AwsDeviceV2:        "virtual",
		AwsLogicalDeviceID: "virtual",
		BGPPeerID:          req.NewBGPPeer.CustomerAddress,
		BGPPeerState:       "available",
		BGPStatus:          "up",
		CustomerAddress:    req.NewBGPPeer.CustomerAddress,
	})

	// Update the Virtual Interface in the database
	err = vifDB.SetVal(req.VirtualInterfaceID, privateVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	returnOk(w, privateVIF)

}

func CreateConnection(w http.ResponseWriter, r *http.Request) {
	dx, err := d.CreateConnection(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	connectionDB, err := db.NewAdapter(dbNameConnection)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()
	// Serialize the struct

	err = connectionDB.SetVal(dx.ConnectionId, &dx)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	// Update the tag database
	resourceTag := d.ResourceTag{
		ResourceArn: d.CreateARN("us-east-1", dx.ConnectionId),
		Tags:        dx.Tags,
	}

	tagDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer tagDB.CloseDbConnection()

	err = tagDB.SetVal(dx.ConnectionId, resourceTag)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	// Return a response
	returnOk(w, dx)
}

// CreateDXGateway creates a Direct Connect Gateway.
// Updates the DB.
func CreateDXGateway(w http.ResponseWriter, r *http.Request) {
	g, err := d.CreateDXGateway(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dxgwDB, err := db.NewAdapter(dbNameDXGwy)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer dxgwDB.CloseDbConnection()

	err = dxgwDB.SetVal(g.DirectConnectGatewayId, g)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	response := struct {
		DirectConnectGateway d.DXGateway `json:"directConnectGateway"`
	}{
		DirectConnectGateway: g,
	}

	returnOk(w, response)
}

// CreatePrivateVirtualInterface
func CreatePrivateVirtualInterface(w http.ResponseWriter, r *http.Request) {
	vif, err := d.CreatePrivateVirtualInterface(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vifDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer vifDB.CloseDbConnection()

	err = vifDB.SetVal(vif.VirtualInterfaceID, vif)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	// Update the tag database
	resourceTag := d.ResourceTag{
		ResourceArn: d.CreateARN("us-east-1", vif.VirtualInterfaceID),
		Tags:        vif.Tags,
	}

	tagDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer tagDB.CloseDbConnection()

	err = tagDB.SetVal(vif.VirtualInterfaceID, resourceTag)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	returnOk(w, vif)
}

// CreatePublicVirtualInterface
func CreatePublicVirtualInterface(w http.ResponseWriter, r *http.Request) {
	vif, err := d.CreatePublicVirtualInterface(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vifDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer vifDB.CloseDbConnection()

	err = vifDB.SetVal(vif.VirtualInterfaceID, vif)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	// Update the tag database
	resourceTag := d.ResourceTag{
		ResourceArn: d.CreateARN("us-east-1", vif.VirtualInterfaceID),
		Tags:        vif.Tags,
	}

	tagDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer tagDB.CloseDbConnection()

	err = tagDB.SetVal(vif.VirtualInterfaceID, resourceTag)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	returnOk(w, vif)
}

// A private virtual interface can be connected to either a Direct Connect gateway or a Virtual Private Gateway (VGW).
func DescribeConnections(w http.ResponseWriter, r *http.Request) {
	request, err := d.DescribeConnections(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	connectionDB, err := db.NewAdapter(dbNameConnection)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	response := d.DescribeConnectionsResponse{
		Connections: []d.Connection{},
	}

	// Find the Connection in the database
	var dx d.Connection
	err = connectionDB.GetVal(request.ConnectionId, &dx)
	if err != nil {
		log.Println("Error in getting connection ID from database", err)
		json.NewEncoder(w).Encode(response)
	}

	response.Connections = append(response.Connections, dx)
	json.NewEncoder(w).Encode(response)
}

// DescribeVirtualInterfaces
func DescribeVirtualInterfaces(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VirtualInterfaceID string `json:"virtualInterfaceId"`
		ConnectionID       string `json:"connectionId"`
	}
	err := d.RequestToJson(r, &req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vifDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer vifDB.CloseDbConnection()

	var privateVIF d.PrivateVirtualInterface

	err = vifDB.GetVal(req.VirtualInterfaceID, &privateVIF)
	if err != nil {
		log.Println("Error in getting virtual interface ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var response struct {
		VirtualInterfaces []d.PrivateVirtualInterface `json:"virtualInterfaces"`
	}

	response.VirtualInterfaces = append(response.VirtualInterfaces, privateVIF)
	json.NewEncoder(w).Encode(response)
}

func DescribeTags(w http.ResponseWriter, r *http.Request) {
	request, err := d.DescribeTags(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	connectionDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	response := d.DescribeTagsResponse{}

	for _, resourceARN := range request.ResourceArns {
		key, err := d.GetIDFromARN(resourceARN)
		if err != nil {
			log.Println("Error in getting key from ARN", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		resourceTag := d.ResourceTag{}
		err = connectionDB.GetVal(key, &resourceTag)
		if err != nil {
			log.Println("Value not found in database", err)
			continue
		}

		response.ResourceTags = append(response.ResourceTags, resourceTag)
	}

	json.NewEncoder(w).Encode(response)
}

// DescribeDXGateways returns a list of Direct Connect Gateways.
// DXGWYs in deleted state are not returned.
func DescribeDXGateways(w http.ResponseWriter, r *http.Request) {
	request, err := d.DescribeDXGateways(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dxgwDB, err := db.NewAdapter(dbNameDXGwy)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer dxgwDB.CloseDbConnection()

	response := d.DescribeDXGatewaysResponse{
		DirectConnectGateways: []d.DXGateway{},
	}

	var g d.DXGateway
	err = dxgwDB.GetVal(request.DirectConnectGatewayId, &g)
	if err != nil {
		log.Println("Error in getting connection ID from database", err)
	}
	if g.DirectConnectGatewayState != "deleted" {
		response.DirectConnectGateways = append(response.DirectConnectGateways, g)
	}

	returnOk(w, response)
}

// DeleteBGPPeer deletes a BGP Peer.
// Check if the Virtual Interface exists in the database.
// Check if the BGP Peer exists in the Virtual Interface.
// Change the state of the BGP Peer to deleted and the status to down.
func DeleteBGPPeer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ASN                int    `json:"asn"`
		BGPPeerID          string `json:"bgpPeerId"`
		CustomerAddress    string `json:"customerAddress"`
		VirtualInterfaceID string `json:"virtualInterfaceId"`
	}
	err := d.RequestToJson(r, &req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vifDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer vifDB.CloseDbConnection()

	var privateVIF d.PrivateVirtualInterface

	err = vifDB.GetVal(req.VirtualInterfaceID, &privateVIF)
	if err != nil {
		log.Println("Error in getting virtual interface ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// check if the BGP Peer already exists
	for key, bgpPeer := range privateVIF.BGPPeers {
		if bgpPeer.CustomerAddress == req.CustomerAddress {
			// Update the BGP Peer
			privateVIF.BGPPeers[key].BGPPeerState = "deleted"
			privateVIF.BGPPeers[key].BGPStatus = "down"
		}
	}

	err = vifDB.SetVal(req.VirtualInterfaceID, privateVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	returnOk(w, privateVIF)
}

// DeleteConnections deletes a connection.
func DeleteConnections(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConnectionID string `json:"connectionId"`
	}
	err := d.RequestToJson(r, &req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	connectionDB, err := db.NewAdapter(dbNameConnection)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	var dx d.Connection
	err = connectionDB.GetVal(req.ConnectionID, &dx)
	if err != nil {
		log.Println("Error in getting connection ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	dx.ConnectionState = "deleted"

	// Delete the connection from the database
	err = connectionDB.SetVal(req.ConnectionID, dx)
	if err != nil {
		log.Println("Error in deleting connection ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	returnOk(w, dx)
}

// DeleteDXGateway deletes a Direct Connect Gateway.
func DeleteDXGateway(w http.ResponseWriter, r *http.Request) {
	request, err := d.DeleteDXGateway(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dxgwDB, err := db.NewAdapter(dbNameDXGwy)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer dxgwDB.CloseDbConnection()

	var g d.DXGateway
	err = dxgwDB.GetVal(request.DirectConnectGatewayId, &g)
	if err != nil {
		log.Println("Error in getting connection ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	g.DirectConnectGatewayState = "deleted"

	err = dxgwDB.SetVal(request.DirectConnectGatewayId, g)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteVirtualInterface deletes a virtual interface.
func DeleteVirtualInterface(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VirtualInterfaceID string `json:"virtualInterfaceId"`
	}
	err := d.RequestToJson(r, &req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	vifDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer vifDB.CloseDbConnection()

	var privateVIF d.PrivateVirtualInterface
	err = vifDB.GetVal(req.VirtualInterfaceID, &privateVIF)
	if err != nil {
		log.Println("Error in getting virtual interface ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	privateVIF.VirtualInterfaceState = "deleted"

	err = vifDB.SetVal(req.VirtualInterfaceID, privateVIF)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TagResource tags a resource
func TagResource(w http.ResponseWriter, r *http.Request) {
	resourceTag, err := d.TagResource(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	connectionDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	key, err := d.GetIDFromARN(resourceTag.ResourceArn)
	if err != nil {
		log.Println("Error in getting key from ARN", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	err = connectionDB.SetVal(key, resourceTag)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateDXGateway updates a Direct Connect Gateway.
// Updates the DB.
func UpdateDXGateway(w http.ResponseWriter, r *http.Request) {
	request, err := d.UpdateDXGateway(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	dxgwDB, err := db.NewAdapter(dbNameDXGwy)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer dxgwDB.CloseDbConnection()

	var g d.DXGateway
	err = dxgwDB.GetVal(request.DirectConnectGatewayId, &g)
	if err != nil {
		log.Println("Error in getting connection ID from database", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	g.DirectConnectGatewayName = request.NewDirectConnectGatewayName

	err = dxgwDB.SetVal(request.DirectConnectGatewayId, g)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	response := struct {
		DirectConnectGateway d.DXGateway `json:"directConnectGateway"`
	}{
		DirectConnectGateway: g,
	}

	returnOk(w, response)
}
