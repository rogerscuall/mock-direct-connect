/*
The file **vif.go** has the implementation of the api calls for private and public virtual interfaces.
It embeds the NewPrivateVirtualInterface and NewPublicVirtualInterface in the struct PrivateVirtualInterface and PublicVirtualInterface respectively.
*/
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
	NewPrivateVirtualInterface
	AmazonSideASN         int                 `json:"amazonSideAsn"`
	AwsDeviceV2           string              `json:"awsDeviceV2"`
	AwsLogicalDeviceID    string              `json:"awsLogicalDeviceId"`
	BGPPeers              []BGPConfig         `json:"bgpPeers"`
	ConnectionID          string              `json:"connectionId"`
	JumboFrameCapable     bool                `json:"jumboFrameCapable"`
	Location              string              `json:"location"`
	OwnerAccount          string              `json:"ownerAccount"`
	Region                string              `json:"region"`
	RouteFilterPrefixes   []RouteFilterPrefix `json:"routeFilterPrefixes"`
	SiteLinkEnabled       bool                `json:"siteLinkEnabled"`
	VirtualInterfaceID    string              `json:"virtualInterfaceId"`
	VirtualInterfaceState string              `json:"virtualInterfaceState"`
	VirtualInterfaceType  string              `json:"virtualInterfaceType"`
}

// Public interface

type CreatePublicVirtualInterfaceRequest struct {
	ConnectionID              string                    `json:"connectionId"`
	NewPublicVirtualInterface NewPublicVirtualInterface `json:"NewPublicVirtualInterface"`
}

type NewPublicVirtualInterface struct {
	AddressFamily       string `json:"addressFamily"`
	AmazonAddress       string `json:"amazonAddress"`
	Asn                 int    `json:"asn"`
	AuthKey             string `json:"authKey"`
	CustomerAddress     string `json:"customerAddress"`
	RouteFilterPrefixes []struct {
		Cidr string `json:"cidr"`
	} `json:"routeFilterPrefixes"`
	Tags                 []DirectConnectTag `json:"tags"`
	VirtualInterfaceName string             `json:"virtualInterfaceName"`
	Vlan                 int                `json:"vlan"`
}

type PublicVirtualInterface struct {
	NewPublicVirtualInterface
	AmazonSideAsn          int                `json:"amazonSideAsn"`
	AwsDeviceV2            string             `json:"awsDeviceV2"`
	AwsLogicalDeviceID     string             `json:"awsLogicalDeviceId"`
	BgpPeers               []BGPConfig        `json:"bgpPeers"`
	ConnectionID           string             `json:"connectionId"`
	CustomerRouterConfig   string             `json:"customerRouterConfig"`
	DirectConnectGatewayID string             `json:"directConnectGatewayId"`
	JumboFrameCapable      bool               `json:"jumboFrameCapable"`
	Location               string             `json:"location"`
	Mtu                    int                `json:"mtu"`
	OwnerAccount           string             `json:"ownerAccount"`
	Region                 string             `json:"region"`
	SiteLinkEnabled        bool               `json:"siteLinkEnabled"`
	Tags                   []DirectConnectTag `json:"tags"`
	VirtualGatewayId       string             `json:"virtualGatewayId"`
	VirtualInterfaceID     string             `json:"virtualInterfaceId"`
	VirtualInterfaceState  string             `json:"virtualInterfaceState"`
	VirtualInterfaceType   string             `json:"virtualInterfaceType"`
	Vlan                   int                `json:"vlan"`
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

// CreatePrivateVirtualInterface uses the request to create a PrivateVirtualInterface.
// The interface is available after creation.
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
	//TODO: Fix this hardcoding
	pvif.AmazonSideASN = req.NewPrivateVirtualInterface.ASN
	pvif.AwsDeviceV2 = "virtual"
	pvif.AwsLogicalDeviceID = "virtual"
	pvif.JumboFrameCapable = false
	//pvif.VirtualInterfaceName = req.NewPrivateVirtualInterface.VirtualInterfaceName

	return pvif, nil
}

// Implementing the Marshaler interface
func (p PublicVirtualInterface) MarshalJSON() ([]byte, error) {
	type Alias PublicVirtualInterface
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&p),
	})
}

// Implementing the Unmarshaler interface
func (p *PublicVirtualInterface) UnmarshalJSON(b []byte) error {
	type Alias PublicVirtualInterface
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	return json.Unmarshal(b, &aux)
}

// CreatePublicVirtualInterface uses the request to create a PublicVirtualInterface.
// The interface is available after creation.
func CreatePublicVirtualInterface(r *http.Request) (PublicVirtualInterface, error) {
	var pvif PublicVirtualInterface
	var req CreatePublicVirtualInterfaceRequest
	err := RequestToJson(r, &req)
	if err != nil {
		return pvif, err
	}
	pvif.ConnectionID = req.ConnectionID
	pvif.NewPublicVirtualInterface = req.NewPublicVirtualInterface
	pvif.VirtualInterfaceID = "dxvif-" + randomString(8)
	pvif.VirtualInterfaceState = "available"
	pvif.VirtualInterfaceType = "public"
	pvif.AmazonSideAsn = 64512
	pvif.AwsDeviceV2 = "virtual"
	pvif.AwsLogicalDeviceID = "virtual"
	pvif.JumboFrameCapable = false
	pvif.Vlan = req.NewPublicVirtualInterface.Vlan
	//pvif.VirtualInterfaceName = req.NewPrivateVirtualInterface.VirtualInterfaceName
	return pvif, nil
}
