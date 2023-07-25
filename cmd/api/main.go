package main

import (
	"flag"
	"net"
	"os"
	"sync"

	d "dx-mock/pkg/dx"

	"github.com/osrg/gobgp/v3/pkg/server"
)

var (
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
	//dbNameBgpPeer = "bgpPeers"
	// dbTransitVIF is the name of the DynamoDB table for Transit Virtual Interfaces
	//dbNameTransitVIF = "transitvifs"
	// dbDXGWyAssociation is the name of the DynamoDB table for Direct Connect Gateway Associations
	dbNameDXGWyAssociation = "dxgwyassociations"
)

type config struct {
	port int
	//env  string
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

func main() {
	var err error
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "Port to listen on")
	flag.BoolVar(&createBgpNeighbor, "enable-bgp", true, "Create BGP Neighbor")
	flag.IntVar(&localBgpAsn, "asn", 65001, "Local BGP ASN")
	flag.StringVar(&logMinLevel, "log", "INFO", "LOG LEVEL (DEBUG|INFO|WARNING|ERROR)")
	
	flag.Parse()

	logLevel := NewCustomLogger(logMinLevel)

	a := &application{
		config:      cfg,
		createBGP:   createBgpNeighbor,
		logger:      logLevel,
		localBgpAsn: localBgpAsn,
	}

	a.logger.Info("create BGP Service is ", createBgpNeighbor)
	if a.createBGP {
		vifs, err := GetActiveVirtualInterfaces()
		if err != nil {
			a.logger.Error("error in getting active Virtual Interfaces", err)
			os.Exit(1)
		}
		vifs = GetVirtualInterfaceWithBgpPeers(vifs)
		err = a.initBgpConfig(vifs)
		if err != nil {
			a.logger.Error("error in initializing BGP configuration", err)
			os.Exit(1)
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
