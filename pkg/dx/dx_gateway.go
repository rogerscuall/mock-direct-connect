package dx

import (
	"net/http"
	"strings"
)

type CreateDXGatewayRequest struct {
	AmazonSideASN            int    `json:"amazonSideAsn"`
	DirectConnectGatewayName string `json:"directConnectGatewayName"`
}

type DXGateway struct {
	CreateDXGatewayRequest
	DirectConnectGatewayId    string `json:"directConnectGatewayId"`
	DirectConnectGatewayState string `json:"directConnectGatewayState"`
	OwnerAccount              string `json:"ownerAccount"`
	StateChangeError          string `json:"stateChangeError"`
}

// DescribeDXGatewaysRequest is the request body for DescribeDXGateways
type DescribeDXGatewaysRequest struct {
	DirectConnectGatewayId string `json:"directConnectGatewayId"`
	MaxResults             int    `json:"maxResults"`
	NextToken              string `json:"nextToken"`
}

// DescribeDXGatewaysResponse is the response body for DescribeDXGateways
type DescribeDXGatewaysResponse struct {
	DirectConnectGateways []DXGateway `json:"directConnectGateways"`
	NextToken             string      `json:"nextToken"`
}

// UpdateDXGatewayRequest is the request body for UpdateDXGateway
type UpdateDXGatewayRequest struct {
	DirectConnectGatewayId      string `json:"directConnectGatewayId"`
	NewDirectConnectGatewayName string `json:"newDirectConnectGatewayName"`
}

// DeleteDXGatewayRequest is the request body for DeleteDXGateway
type DeleteDXGatewayRequest struct {
	DirectConnectGatewayId string `json:"directConnectGatewayId"`
}

// CreateDXGateway creates a new DXGateway from the request body
func CreateDXGateway(r *http.Request) (DXGateway, error) {
	var dxg DXGateway
	err := RequestToJson(r, &dxg)
	if err != nil {
		return dxg, err
	}

	dxg.DirectConnectGatewayId = CreateDxGatewayID()
	dxg.DirectConnectGatewayState = "available"
	// TODO: This should not be hardcoded
	dxg.OwnerAccount = "123456789012"

	return dxg, nil
}

// DescribeConnections
func DescribeDXGateways(r *http.Request) (DescribeDXGatewaysRequest, error) {
	var dxg DescribeDXGatewaysRequest
	err := RequestToJson(r, &dxg)
	if err != nil {
		return dxg, err
	}
	return dxg, nil
}

// UpdateDXGateway
func UpdateDXGateway(r *http.Request) (UpdateDXGatewayRequest, error) {
	var dxg UpdateDXGatewayRequest
	err := RequestToJson(r, &dxg)
	if err != nil {
		return dxg, err
	}
	return dxg, nil
}

// DeleteDXGateway
func DeleteDXGateway(r *http.Request) (DeleteDXGatewayRequest, error) {
	var dxg DeleteDXGatewayRequest
	err := RequestToJson(r, &dxg)
	if err != nil {
		return dxg, err
	}
	return dxg, nil
}

// CreateDxGatewayID creates a random string representing the DXGateway ID
func CreateDxGatewayID() string {
	r1 := randomString(8)
	r2 := randomString(4)
	r3 := randomString(4)
	r4 := randomString(4)
	r5 := randomString(12)
	return strings.ToLower(r1 + "-" + r2 + "-" + r3 + "-" + r4 + "-" + r5)
}
