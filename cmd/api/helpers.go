package main

import (
	"dx-mock/adapters/db"
	"dx-mock/pkg/bgp"
	d "dx-mock/pkg/dx"
	"encoding/json"
	"net"
	"net/http"
)

func returnOk(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func GetVirtualInterfaces() ([]d.PrivateVirtualInterface, error) {
	var vifs []d.PrivateVirtualInterface
	connectionDB, err := db.NewAdapter(dbNameVIF)
	if err != nil {
		return vifs, err
	}
	defer connectionDB.CloseDbConnection()
	keys, err := connectionDB.GetKeys()
	if err != nil {
		return vifs, err
	}
	for _, key := range keys {
		var vif d.PrivateVirtualInterface
		err := connectionDB.GetVal(string(key), &vif)
		if err != nil {
			return vifs, err
		}
		vifs = append(vifs, vif)
	}
	return vifs, nil
}

func GetActiveVirtualInterfaces() ([]d.PrivateVirtualInterface, error) {
	vifs, err := GetVirtualInterfaces()
	if err != nil {
		return vifs, err
	}
	var activeVifs []d.PrivateVirtualInterface
	for _, vif := range vifs {
		if vif.VirtualInterfaceState == "available" {
			activeVifs = append(activeVifs, vif)
		}
	}
	return activeVifs, nil
}

func GetDirectConnectGatewaysAssociations() ([]d.DirectConnectGatewayAssociationResponse, error) {
	var dx []d.DirectConnectGatewayAssociationResponse
	connectionDB, err := db.NewAdapter(dbNameDXGWyAssociation)
	if err != nil {
		return dx, err
	}
	defer connectionDB.CloseDbConnection()
	keys, err := connectionDB.GetKeys()
	if err != nil {
		return dx, err
	}
	for _, key := range keys {
		var dxa d.DirectConnectGatewayAssociationResponse
		err := connectionDB.GetVal(string(key), &dxa)
		if err != nil {
			return dx, err
		}
		dx = append(dx, dxa)
	}
	return dx, nil
}

// GetAssociatedDirectConnectGateways returns a list of Direct Connect Gateways associated with the Virtual Interface.
// This will return a list of Direct Connect Gateways.
func GetAssociatedDirectConnectGateways() ([]d.DirectConnectGatewayAssociationResponse, error) {
	dx, err := GetDirectConnectGatewaysAssociations()
	if err != nil {
		return dx, err
	}
	var dxa []d.DirectConnectGatewayAssociationResponse
	for _, dx := range dx {
		if dx.DirectConnectGatewayAssociation.AssociationState == "associated" {
			dxa = append(dxa, dx)
		}
	}
	return dxa, nil
}

// GetVirtualInterfaceWithBgpPeers returns a list of Virtual Interfaces with BGP Peers.
// This will return a list of Virtual Interfaces with BGP Peers.
// If the Virtual Interface field BGPPeers is not empty it will be returned here.
func GetVirtualInterfaceWithBgpPeers(pvifs []d.PrivateVirtualInterface) []d.PrivateVirtualInterface {
	var vifs []d.PrivateVirtualInterface
	for _, pvif := range pvifs {
		if len(pvif.BGPPeers) > 0 {
			vifs = append(vifs, pvif)
		}
	}
	return vifs
}

// initBgpConfig initializes the BGP configuration for the Virtual Interface.
// For all active Virtual Interfaces it will create a BGP server and BGP peers if they are not deleted.
// For all the DX Associations that are not "disassociated" it will add the path allowed prefixes to the BGP server.
func (a *application) initBgpConfig(vifs []d.PrivateVirtualInterface) error {
	var err error
	a.logger.Info("creating BGP service")
	a.primaryIP, err = bgp.GetPrimaryIP()
	if err != nil {
		a.logger.Error("error in getting primary IP address ", err)
		return err
	}
	a.logger.Info("primary IP address is ", a.primaryIP)
	a.serverBgp, err = bgp.CreateBgpServer(a.localBgpAsn, a.primaryIP)
	if err != nil {
		a.logger.Error("error in creating BGP server ", err)
		return err
	}
	a.logger.Info("creating BGP peers")
	for _, vif := range vifs {
		a.logger.Info("virtual Interface:", vif.VirtualInterfaceID)
		for _, bgpPeer := range vif.BGPPeers {
			// Check if the BGP peer is active
			if bgpPeer.BGPPeerState == "deleted" {
				continue
			}
			a.logger.Info("BGP Peer:", bgpPeer.BGPPeerID, " on port TCP 179")
			err = bgp.CreateBGPPeer(a.serverBgp, bgpPeer.ASN, net.ParseIP(bgpPeer.CustomerAddress))
			if err != nil {
				a.logger.Error("error in creating BGP peer", err)
				return err
			}
			a.bgpPeers = append(a.bgpPeers, bgpPeer)
		}
	}
	a.logger.Info("adding paths to BGP server")
	dx, err := GetAssociatedDirectConnectGateways()
	if err != nil {
		a.logger.Error("error in getting associated Direct Connect Gateways", err)
		return err
	}
	for _, dx := range dx {
		for _, prefix := range dx.DirectConnectGatewayAssociation.AllowedPrefixesToDirectConnectGateway {
			a.logger.Info("adding path to BGP server for prefix:", prefix.Cidr)
			_, path, err := net.ParseCIDR(prefix.Cidr)
			if err != nil {
				a.logger.Error("error in parsing CIDR", err)
				return err
			}
			err = bgp.AddPath(a.serverBgp, *path, a.primaryIP)
			if err != nil {
				a.logger.Error("error in adding path to BGP server", err)
				return err
			}
		}
	}
	return nil
}
