package dx

import (
	"encoding/json"
	"net/http"
)

// CreateConnectionRequest is the request body for CreateConnection
type CreateConnectionRequest struct {
	Bandwidth      string `json:"bandwidth"`
	ConnectionName string `json:"connectionName"`
	LagId          string `json:"lagId"`
	Location       string `json:"location"`
	ProviderName   string `json:"providerName"`
	RequestMACSec  bool   `json:"requestMACSec"`
	Action         string `json:"Action"`
}

// CreateConnectionResponse is the response body for CreateConnection
type CreateConnectionResponse struct {
	AwsDevice            string             `json:"awsDevice"`
	AwsDeviceV2          string             `json:"awsDeviceV2"`
	AwsLogicalDeviceId   string             `json:"awsLogicalDeviceId"`
	Bandwidth            string             `json:"bandwidth"`
	ConnectionId         string             `json:"connectionId"`
	ConnectionName       string             `json:"connectionName"`
	ConnectionState      string             `json:"connectionState"`
	EncryptionMode       string             `json:"encryptionMode"`
	HasLogicalRedundancy string             `json:"hasLogicalRedundancy"`
	JumboFrameCapable    bool               `json:"jumboFrameCapable"`
	LagId                string             `json:"lagId"`
	LoaIssueTime         int64              `json:"loaIssueTime"`
	Location             string             `json:"location"`
	MacSecCapable        bool               `json:"macSecCapable"`
	MacSecKeys           []DirectConnectKey `json:"macSecKeys"`
	OwnerAccount         string             `json:"ownerAccount"`
	PartnerName          string             `json:"partnerName"`
	PortEncryptionStatus string             `json:"portEncryptionStatus"`
	ProviderName         string             `json:"providerName"`
	Region               string             `json:"region"`
	Tags                 []DirectConnectTag `json:"tags"`
	Vlan                 int                `json:"vlan"`
}

// DirectConnectKey is the key for MACSec used by CreateConnectionResponse
type DirectConnectKey struct {
	Ckn       string `json:"ckn"`
	SecretArn string `json:"secretARN"`
	StartOn   string `json:"startOn"`
	State     string `json:"state"`
}

// DescribeConnectionsRequest is the request body for DescribeConnections
type DescribeConnectionsRequest struct {
	ConnectionId string `json:"connectionId"`
}

// DescribeConnectionsResponse is the response body for DescribeConnections
type DescribeConnectionsResponse struct {
	Connections []CreateConnectionResponse `json:"connections"`
}

func NewCreateConnectionResponse(dx CreateConnectionRequest) *CreateConnectionResponse {
	return &CreateConnectionResponse{
		AwsDevice:            "123",
		AwsDeviceV2:          "123",
		AwsLogicalDeviceId:   "123",
		Bandwidth:            dx.Bandwidth,
		ConnectionId:         "dxcon-1234567890",
		ConnectionName:       dx.ConnectionName,
		ConnectionState:      "available",
		EncryptionMode:       "no_encrypt",
		HasLogicalRedundancy: "yes",
		JumboFrameCapable:    false,
		LagId:                dx.LagId,
		LoaIssueTime:         0,
		Location:             dx.Location,
		MacSecCapable:        false,
		MacSecKeys:           nil,
		OwnerAccount:         "123",
		PartnerName:          "p1",
		PortEncryptionStatus: "Encryption Down",
		ProviderName:         "p1",
		Region:               "us-east-1",
		Tags:                 nil,
		Vlan:                 123,
	}
}

func CreateConnection(r *http.Request) (*CreateConnectionResponse, error) {
	var dx CreateConnectionRequest

	// Unmarshal the body
	err := RequestToJson(r, &dx)
	if err != nil {
		return nil, err
	}

	// Return a response
	return NewCreateConnectionResponse(dx), nil
}

func RequestToJson(r *http.Request, v interface{}) error {
	// Unmarshal the body
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}

	return nil
}

func DescribeConnections(r *http.Request, dx *CreateConnectionResponse) (*DescribeConnectionsResponse, error) {
	var request DescribeConnectionsRequest

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return nil, err
	}

	response := DescribeConnectionsResponse{
		Connections: []CreateConnectionResponse{*dx},
	}

	return &response, nil
}
