package dx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CreatePrivateVirtualInterfaceRequest is the request body for CreatePrivateVirtualInterface
type CreatePrivateVirtualInterfaceRequest struct {
	ConnectionID               string                     `json:"connectionId"`
	NewPrivateVirtualInterface NewPrivateVirtualInterface `json:"newPrivateVirtualInterface"`
}

type NewPrivateVirtualInterface struct {
	AddressFamily          string             `json:"addressFamily"`
	AmazonAddress          string             `json:"amazonAddress"`
	ASN                    int                `json:"asn"`
	AuthKey                string             `json:"authKey"`
	CustomerAddress        string             `json:"customerAddress"`
	DirectConnectGatewayID string             `json:"directConnectGatewayId"`
	EnableSiteLink         bool               `json:"enableSiteLink"`
	MTU                    int                `json:"mtu"`
	Tags                   []DirectConnectTag `json:"tags"`
	VirtualGatewayID       string             `json:"virtualGatewayId"`
	VirtualInterfaceName   string             `json:"virtualInterfaceName"`
	VLAN                   int                `json:"vlan"`
}

type BGPConfig struct {
	AddressFamily      string `json:"addressFamily"`
	AmazonAddress      string `json:"amazonAddress"`
	ASN                int    `json:"asn"`
	AuthKey            string `json:"authKey"`
	AwsDeviceV2        string `json:"awsDeviceV2"`
	AwsLogicalDeviceID string `json:"awsLogicalDeviceId"`
	BGPPeerID          string `json:"bgpPeerId"`
	BGPPeerState       string `json:"bgpPeerState"`
	BGPStatus          string `json:"bgpStatus"`
	CustomerAddress    string `json:"customerAddress"`
}

type RouteFilterPrefix struct {
	CIDR string `json:"cidr"`
}

// PrivateVirtualInterface is the response body for CreatePrivateVirtualInterface
// It is equivalent to Response of the action CreatePrivateVirtualInterface.
// https://docs.aws.amazon.com/directconnect/latest/APIReference/API_CreatePrivateVirtualInterface.html
type PrivateVirtualInterface struct {
	// Embedded NewPrivateVirtualInterface
	ConnectionID          string              `json:"connectionId"`
	AmazonSideASN         int                 `json:"amazonSideAsn"`
	AwsDeviceV2           string              `json:"awsDeviceV2"`
	AwsLogicalDeviceID    string              `json:"awsLogicalDeviceId"`
	BGPPeers              []BGPConfig         `json:"bgpPeers"`
	JumboFrameCapable     bool                `json:"jumboFrameCapable"`
	Location              string              `json:"location"`
	OwnerAccount          string              `json:"ownerAccount"`
	Region                string              `json:"region"`
	RouteFilterPrefixes   []RouteFilterPrefix `json:"routeFilterPrefixes"`
	SiteLinkEnabled       bool                `json:"siteLinkEnabled"`
	VirtualInterfaceID    string              `json:"virtualInterfaceId"`
	VirtualInterfaceState string              `json:"virtualInterfaceState"`
	VirtualInterfaceType  string              `json:"virtualInterfaceType"`
	NewPrivateVirtualInterface
}

// Implementing the Marshaler interface
func (p PrivateVirtualInterface) MarshalJSON() ([]byte, error) {
	type Alias PrivateVirtualInterface
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&p),
	})
}

// Implementing the Unmarshaler interface
func (p *PrivateVirtualInterface) UnmarshalJSON(b []byte) error {
	type Alias PrivateVirtualInterface
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	return json.Unmarshal(b, &aux)
}

// CreatePrivateVirtualInterface
func CreatePrivateVirtualInterface(r *http.Request) (PrivateVirtualInterface, error) {
	var pvif PrivateVirtualInterface
	var req CreatePrivateVirtualInterfaceRequest
	err := RequestToJson(r, &req)
	if err != nil {
		return pvif, err
	}
	if req.NewPrivateVirtualInterface.VirtualGatewayID != "" && req.NewPrivateVirtualInterface.DirectConnectGatewayID != "" {
		return pvif, fmt.Errorf("VirtualGatewayID and DirectConnectGatewayID cannot be specified together")
	}
	pvif.ConnectionID = req.ConnectionID
	pvif.NewPrivateVirtualInterface = req.NewPrivateVirtualInterface
	pvif.VirtualInterfaceID = "dxvif-" + randomString(8)
	pvif.VirtualInterfaceState = "available"
	pvif.VirtualInterfaceType = "private"
	pvif.AmazonSideASN = 64512
	pvif.AwsDeviceV2 = "virtual"
	pvif.AwsLogicalDeviceID = "virtual"
	pvif.JumboFrameCapable = false
	//pvif.VirtualInterfaceName = req.NewPrivateVirtualInterface.VirtualInterfaceName

	return pvif, nil
}
