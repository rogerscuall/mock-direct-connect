package main

import (
	d "dx-mock/pkg/dx"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var (
	dx d.Connection
)

func (a *application) handleRequest(w http.ResponseWriter, r *http.Request) {
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
	log.Println("Request for:", action)
	switch action {
	case "CreateBGPPeer":
		a.CreateBGPPeer(w, r)
	case "CreateConnection":
		a.CreateConnection(w, r)
	case "CreateDirectConnectGateway":
		a.CreateDXGateway(w, r)
	case "CreateDirectConnectGatewayAssociation":
		a.CreateDirectConnectGatewayAssociation(w, r)
	case "CreatePrivateVirtualInterface":
		a.CreatePrivateVirtualInterface(w, r)
	case "CreatePublicVirtualInterface":
		a.CreatePublicVirtualInterface(w, r)
	case "CreateTransitVirtualInterface":
		a.CreateTransitVirtualInterface(w, r)
	case "DeleteBGPPeer":
		a.DeleteBGPPeer(w, r)
	case "DeleteConnection":
		a.DeleteConnections(w, r)
	case "DeleteDirectConnectGateway":
		a.DeleteDXGateway(w, r)
	case "DeleteDirectConnectGatewayAssociation":
		a.DeleteDirectConnectGatewayAssociation(w, r)
	case "DeleteVirtualInterface":
		a.DeleteVirtualInterface(w, r)
	case "DescribeConnections":
		a.DescribeConnections(w, r)
	case "DescribeDirectConnectGateways":
		a.DescribeDXGateways(w, r)
	case "DescribeDirectConnectGatewayAssociations":
		a.DescribeDirectConnectGatewayAssociations(w, r)
	case "DescribeVirtualInterfaces":
		a.DescribeVirtualInterfaces(w, r)
	case "DescribeTags":
		a.DescribeTags(w, r)
	case "TagResource":
		a.TagResource(w, r)
	case "UpdateConnection":
		err := d.UpdateConnection(r, &dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(dx)

		return
	case "UpdateDirectConnectGateway":
		a.UpdateDXGateway(w, r)
	}
}
