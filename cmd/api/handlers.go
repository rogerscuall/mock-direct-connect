package main

import (
	"dx-mock/adapters/db"
	d "dx-mock/pkg/dx"
	"encoding/json"
	"log"
	"net/http"
)

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
