package dx

import (
	"net/http"
)

type DirectConnectGatewayPrefix struct {
	Cidr string `json:"cidr"`
}

// CreateDirectConnectGatewayAssociationRequest is the request body for CreateDirectConnectGatewayAssociation
type CreateDirectConnectGatewayAssociationRequest struct {
	AddAllowedPrefixesToDirectConnectGateway []DirectConnectGatewayPrefix `json:"addAllowedPrefixesToDirectConnectGateway"`
	DirectConnectGatewayId                   string                       `json:"directConnectGatewayId"`
	GatewayId                                string                       `json:"gatewayId"`
	VirtualGatewayId                         string                       `json:"virtualGatewayId"`
}

type AssociatedGateway struct {
	Id           string `json:"id"`
	OwnerAccount string `json:"ownerAccount"`
	Region       string `json:"region"`
	Type         string `json:"type"`
}

type DirectConnectGatewayAssociation struct {
	AllowedPrefixesToDirectConnectGateway []DirectConnectGatewayPrefix `json:"allowedPrefixesToDirectConnectGateway"`
	AssociatedGateway                     AssociatedGateway            `json:"associatedGateway"`
	AssociationId                         string                       `json:"associationId"`
	AssociationState                      string                       `json:"associationState"`
	DirectConnectGatewayId                string                       `json:"directConnectGatewayId"`
	DirectConnectGatewayOwnerAccount      string                       `json:"directConnectGatewayOwnerAccount"`
	StateChangeError                      string                       `json:"stateChangeError"`
	VirtualGatewayId                      string                       `json:"virtualGatewayId"`
	VirtualGatewayOwnerAccount            string                       `json:"virtualGatewayOwnerAccount"`
	VirtualGatewayRegion                  string                       `json:"virtualGatewayRegion"`
}

// DirectConnectGatewayAssociationResponse is the response body for CreateDirectConnectGatewayAssociation
type DirectConnectGatewayAssociationResponse struct {
	DirectConnectGatewayAssociation DirectConnectGatewayAssociation `json:"directConnectGatewayAssociation"`
}

func CreateDirectConnectGatewayAssociation(r *http.Request) (DirectConnectGatewayAssociationResponse, error) {
	var dxga DirectConnectGatewayAssociation
	var dxgar DirectConnectGatewayAssociationResponse
	var req CreateDirectConnectGatewayAssociationRequest
	err := RequestToJson(r, &req)
	if err != nil {
		return dxgar, err
	}
	dxga.AllowedPrefixesToDirectConnectGateway = req.AddAllowedPrefixesToDirectConnectGateway
	dxga.AssociatedGateway.Id = req.GatewayId
	//IMPROVE: This is a mock implementation. The real implementation dont matter.
	dxga.AssociatedGateway.OwnerAccount = "123456789012"
	//IMPROVE: This is a mock implementation. The real implementation dont matter.
	dxga.AssociatedGateway.Region = "us-east-1"
	//FIXME: This is hardcoded. At the current moment only TGW is needed.
	dxga.AssociatedGateway.Type = "transitGateway"
	dxga.AssociationId = createDxGatewayAssociationID()
	dxga.AssociationState = "associated"
	dxga.DirectConnectGatewayId = req.DirectConnectGatewayId
	//IMPROVE: This is a mock implementation. The real implementation dont matter.
	dxga.DirectConnectGatewayOwnerAccount = "123456789012"
	dxga.StateChangeError = ""
	dxga.VirtualGatewayId = req.VirtualGatewayId
	//IMPROVE: This is a mock implementation. The real implementation dont matter.
	dxga.VirtualGatewayOwnerAccount = "123456789012"
	//IMPROVE: This is a mock implementation. The real implementation dont matter.
	dxga.VirtualGatewayRegion = "us-east-1"
	dxgar.DirectConnectGatewayAssociation = dxga
	return dxgar, nil
}
