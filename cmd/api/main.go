package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"

	"dx-mock/pkg/bgp"
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

	a.logger.Info("the value of createBgpNeighbor is: ", createBgpNeighbor)
	if a.createBGP {
		a.logger.Info("creating BGP service")
		a.primaryIP, err = bgp.GetPrimaryIP()
		if err != nil {
			a.logger.Error("error in getting primary IP address", err)
			os.Exit(1)
		}
		a.logger.Info("primary IP address is", a.primaryIP)
		a.serverBgp, err = bgp.CreateBgpServer(a.localBgpAsn, a.primaryIP)
		if err != nil {
			a.logger.Error("error in creating BGP server", err)
			os.Exit(1)
		}
		vifs, err := GetActiveVirtualInterfaces()
		if err != nil {
			a.logger.Error("error in getting active Virtual Interfaces", err)
			os.Exit(1)
		}
		vifs = GetVirtualInterfaceWithBgpPeers(vifs)
		log.Println("Creating BGP peers")
		for _, vif := range vifs {
			a.logger.Info("virtual Interface:", vif.VirtualInterfaceID)
			for _, bgpPeer := range vif.BGPPeers {
				a.logger.Info("BGP Peer:", bgpPeer.BGPPeerID)
				err = bgp.CreateBGPPeer(a.serverBgp, bgpPeer.ASN, net.ParseIP(bgpPeer.CustomerAddress))
				if err != nil {
					a.logger.Error("Error in creating BGP peer", err)
					os.Exit(1)
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
