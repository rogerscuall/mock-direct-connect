package dx

type CreateBgpPeerRequest struct {
	NewBGPPeer         NewBGPPeer `json:"newBGPPeer"`
	VirtualInterfaceID string `json:"virtualInterfaceId"`
}

type NewBGPPeer struct {
	AddressFamily   string `json:"addressFamily"`
	AmazonAddress   string `json:"amazonAddress"`
	Asn             int    `json:"asn"`
	AuthKey         string `json:"authKey"`
	CustomerAddress string `json:"customerAddress"`
}

// type BGPPeer struct {
// 	VirtualInterface VirtualInterface `json:"virtualInterface"`
// }

// type VirtualInterface struct {
// 	AddressFamily          string              `json:"addressFamily"`
// 	AmazonAddress          string              `json:"amazonAddress"`
// 	AmazonSideAsn          int                 `json:"amazonSideAsn"`
// 	Asn                    int                 `json:"asn"`
// 	AuthKey                string              `json:"authKey"`
// 	AwsDeviceV2            string              `json:"awsDeviceV2"`
// 	AwsLogicalDeviceId     string              `json:"awsLogicalDeviceId"`
// 	BgpPeers               []BgpPeer           `json:"bgpPeers"`
// 	ConnectionId           string              `json:"connectionId"`
// 	CustomerAddress        string              `json:"customerAddress"`
// 	CustomerRouterConfig   string              `json:"customerRouterConfig"`
// 	DirectConnectGatewayId string              `json:"directConnectGatewayId"`
// 	JumboFrameCapable      bool                `json:"jumboFrameCapable"`
// 	Location               string              `json:"location"`
// 	Mtu                    int                 `json:"mtu"`
// 	OwnerAccount           string              `json:"ownerAccount"`
// 	Region                 string              `json:"region"`
// 	RouteFilterPrefixes    []RouteFilterPrefix `json:"routeFilterPrefixes"`
// 	SiteLinkEnabled        bool                `json:"siteLinkEnabled"`
// 	Tags                   []DirectConnectTag  `json:"tags"`
// 	VirtualGatewayId       string              `json:"virtualGatewayId"`
// 	VirtualInterfaceId     string              `json:"virtualInterfaceId"`
// 	VirtualInterfaceName   string              `json:"virtualInterfaceName"`
// 	VirtualInterfaceState  string              `json:"virtualInterfaceState"`
// 	VirtualInterfaceType   string              `json:"virtualInterfaceType"`
// 	Vlan                   int                 `json:"vlan"`
// }

// type BgpPeer struct {
// 	AddressFamily      string `json:"addressFamily"`
// 	AmazonAddress      string `json:"amazonAddress"`
// 	Asn                int    `json:"asn"`
// 	AuthKey            string `json:"authKey"`
// 	AwsDeviceV2        string `json:"awsDeviceV2"`
// 	AwsLogicalDeviceId string `json:"awsLogicalDeviceId"`
// 	BgpPeerId          string `json:"bgpPeerId"`
// 	BgpPeerState       string `json:"bgpPeerState"`
// 	BgpStatus          string `json:"bgpStatus"`
// 	CustomerAddress    string `json:"customerAddress"`
// }

// // Implementing the Marshaler interface
// func (b BgpPeer) MarshalJSON() ([]byte, error) {
// 	type Alias BgpPeer
// 	return json.Marshal(&struct {
// 		*Alias
// 	}{
// 		Alias: (*Alias)(&b),
// 	})
// }

// // Implementing the Unmarshaler interface
// func (b *BgpPeer) UnmarshalJSON(data []byte) error {
// 	type Alias BgpPeer
// 	aux := &struct {
// 		*Alias
// 	}{
// 		Alias: (*Alias)(b),
// 	}
// 	return json.Unmarshal(data, &aux)
// }

// // CreateBgpPeer creates a BGP Peer
// // Takes a request and return a BGPPeer and an error
// // It is just a helper to parse the request inside the handler for CreateBgpPeer.
// // func CreateBgpPeer(r *http.Request) (BGPPeer, VirtualInterface, error) {
// // 	var bgpPeer BGPPeer
// // 	var req CreateBgpPeerRequest
// // 	err := RequestToJson(r, &req)
// // 	if err != nil {
// // 		return bgpPeer, VirtualInterface{}, err
// // 	}
// // 	bgpPeer.VirtualInterface.AddressFamily = req.NewBGPPeer.AddressFamily
// // 	bgpPeer.VirtualInterface.AmazonAddress = req.NewBGPPeer.AmazonAddress
// // 	bgpPeer.VirtualInterface.AmazonSideAsn = req.NewBGPPeer.Asn
// // 	bgpPeer.VirtualInterface.Asn = req.NewBGPPeer.Asn
// // 	bgpPeer.VirtualInterface.AuthKey = req.NewBGPPeer.AuthKey
// // 	bgpPeer.VirtualInterface.CustomerAddress = req.NewBGPPeer.CustomerAddress
// // 	bgpPeer.VirtualInterface.AwsDeviceV2 = "virtual"
// // 	bgpPeer.VirtualInterface.AwsLogicalDeviceId = "virtual"
// // 	bgpPeer.VirtualInterface.BgpPeers = append(bgpPeer.VirtualInterface.BgpPeers, BgpPeer{
// // 		AddressFamily:      req.NewBGPPeer.AddressFamily,
// // 		AmazonAddress:      req.NewBGPPeer.AmazonAddress,
// // 		Asn:                req.NewBGPPeer.Asn,
// // 		AuthKey:            req.NewBGPPeer.AuthKey,
// // 		AwsDeviceV2:        "virtual",
// // 		AwsLogicalDeviceId: "virtual",
// // 		BgpPeerId:          req.NewBGPPeer.CustomerAddress,
// // 		BgpPeerState:       "available",
// // 		BgpStatus:          "up",
// // 		CustomerAddress:    req.NewBGPPeer.CustomerAddress,
// // 	})
// // 	// Not sure about this one, needs verification.
// // 	bgpPeer.VirtualInterface.ConnectionId = req.VirtualInterfaceId
// // 	return bgpPeer, nil
// // }
