package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"dx-mock/pkg/bgp"
	d "dx-mock/pkg/dx"

	"github.com/osrg/gobgp/v3/pkg/server"
)

var (
	dx                d.Connection
	logMinLevel       string
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
	dbNameVIF = "vifs2"
	// dbBgpPeer is the name of the DynamoDB table for BGP Peers
	dbNameBgpPeer = "bgpPeers"
	// dbTransitVIF is the name of the DynamoDB table for Transit Virtual Interfaces
	dbNameTransitVIF = "transitvifs"
	// dbDXGWyAssociation is the name of the DynamoDB table for Direct Connect Gateway Associations
	dbNameDXGWyAssociation = "dxgwyassociations"
)

type config struct {
	port int
	env  string
}

type application struct {
	config      config
	logger      *CustomLogger
	wg          sync.WaitGroup
	createBGP   bool
	primaryIP   net.IP
	serverBgp   *server.BgpServer
	localBgpAsn int
	bgpPeers    []d.BGPConfig
}

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

func main() {
	var err error
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "Port to listen on")
	flag.StringVar(&logMinLevel, "log", "INFO", "LOG LEVEL (DEBUG|INFO|WARNING|ERROR)")
	flag.IntVar(&localBgpAsn, "asn", 65001, "Local BGP ASN")

	// Load createBgpNeighbor from the environment variable CREATE_BGP_NEIGHBOR
	createBgpNeighbor = os.Getenv("CREATE_BGP_NEIGHBOR") == "true"

	logLevel := NewCustomLogger(logMinLevel)

	a := &application{
		config:      cfg,
		createBGP:   createBgpNeighbor,
		logger:      logLevel,
		localBgpAsn: localBgpAsn,
	}

	a.logger.Info("the value of createBgpNeighbor is:", createBgpNeighbor)
	if a.createBGP {
		a.logger.Info("creating BGP service")
		a.primaryIP, err = bgp.GetPrimaryIP()
		if err != nil {
			a.logger.Error("error in getting primary IP address", err)
			os.Exit(1)
		}
		a.serverBgp, err = bgp.CreateBgpServer(a.localBgpAsn, a.primaryIP)
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
				err = bgp.CreateBGPPeer(a.serverBgp, bgpPeer.ASN, net.ParseIP(bgpPeer.CustomerAddress))
				if err != nil {
					log.Println("Error in creating BGP peer", err)
				}
				a.bgpPeers = append(a.bgpPeers, bgpPeer)
			}
		}
	}
	//http.HandleFunc("/",a.handleRequest)
	err = a.serve()
	if err != nil {
		a.logger.Error(err)
	}
	// fmt.Println("Mock Direct Connect API server listening on port 8080")
	// http.ListenAndServe(":8080", nil)
}
