package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"dx-mock/pkg/bgp"
	d "dx-mock/pkg/dx"
)

var (
	dx                d.Connection
	createBgpNeighbor bool
	localBgpAsn       = 65001
)

const (
	// dbNameConnection is the name of the DynamoDB table for connections
	dbNameConnection = "connection"
	// dbNameDXGwy is the name of the DynamoDB table for Direct Connect Gateways
	dbNameTags = "tags"
	// dbNameDXGwy is the name of the DynamoDB table for Direct Connect Gateways
	dbNameDXGwy = "dxgwys"
	// dbNameVIF is the name of the DynamoDB table for Virtual Interfaces
	dbNameVIF = "vifs"
	// dbBgpPeer is the name of the DynamoDB table for BGP Peers
	dbNameBgpPeer = "bgpPeers"
	// dbTransitVIF is the name of the DynamoDB table for Transit Virtual Interfaces
	dbNameTransitVIF = "transitvifs"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Get the Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-amz-json-1.1" {
		log.Println("Content-Type is not application/x-amz-json-1.1")
	}
	// Get the target
	serviceAction := strings.Split(r.Header.Get("X-Amz-Target"), ".")
	if len(serviceAction) != 2 {
		log.Println("X-Amz-Target is not in the correct format")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	service := serviceAction[0]
	if service != "OvertureService" {
		log.Println("Service is not OvertureService")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	action := serviceAction[1]
	log.Println("Request for: ", action)
	switch action {
	case "CreateBGPPeer":
		CreateBGPPeer(w, r)
	case "CreateConnection":
		CreateConnection(w, r)
	case "CreateDirectConnectGateway":
		CreateDXGateway(w, r)
	case "CreatePrivateVirtualInterface":
		CreatePrivateVirtualInterface(w, r)
	case "CreatePublicVirtualInterface":
		CreatePublicVirtualInterface(w, r)
	case "CreateTransitVirtualInterface":
		CreateTransitVirtualInterface(w, r)
	case "DeleteBGPPeer":
		DeleteBGPPeer(w, r)
	case "DeleteConnection":
		DeleteConnections(w, r)
	case "DeleteDirectConnectGateway":
		DeleteDXGateway(w, r)
	case "DeleteVirtualInterface":
		DeleteVirtualInterface(w, r)
	case "DescribeConnections":
		DescribeConnections(w, r)
	case "DescribeDirectConnectGateways":
		DescribeDXGateways(w, r)
	case "DescribeVirtualInterfaces":
		DescribeVirtualInterfaces(w, r)
	case "DescribeTags":
		DescribeTags(w, r)
	case "TagResource":
		TagResource(w, r)
	case "UpdateConnection":
		err := d.UpdateConnection(r, &dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(dx)

		return
	case "UpdateDirectConnectGateway":
		UpdateDXGateway(w, r)
	}
}

func main() {
	// Load createBgpNeighbor from the environment variable CREATE_BGP_NEIGHBOR
	createBgpNeighbor = os.Getenv("CREATE_BGP_NEIGHBOR") == "true"
	log.Println("The value of createBgpNeighbor is:", createBgpNeighbor)
	if createBgpNeighbor {
		log.Println("Creating BGP service")
		ipAddress, err := bgp.GetPrimaryIP()
		if err != nil {
			log.Panic("Error in getting primary IP address", err)
		}
		serverBgp, err := bgp.CreateBgpServer(localBgpAsn, ipAddress)
		if err != nil {
			log.Panic("Error in creating BGP server", err)
		}
		vifs, err := GetActiveVirtualInterfaces()
		if err != nil {
			log.Panic("Error in getting active Virtual Interfaces", err)
		}
		vifs = GetVirtualInterfaceWithBgpPeers(vifs)
		log.Println("Creating BGP peers")
		for _, vif := range vifs {
			log.Println("Virtual Interface:", vif.VirtualInterfaceID)
			for _, bgpPeer := range vif.BGPPeers {
				log.Println("BGP Peer:", bgpPeer.BGPPeerID)
				err = bgp.CreateBGPPeer(serverBgp, bgpPeer.ASN, net.ParseIP(bgpPeer.CustomerAddress))
				if err != nil {
					log.Println("Error in creating BGP peer", err)
				}
			}
		}
	}
	http.HandleFunc("/", handleRequest)
	fmt.Println("Mock Direct Connect API server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
