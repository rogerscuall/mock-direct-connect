package main

import (
	"dx-mock/adapters/db"
	d "dx-mock/pkg/dx"
	"encoding/json"
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
