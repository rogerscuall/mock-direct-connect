package dx

import (
	"encoding/json"
	"net/http"
)

// CreateConnectionRequest is the request body for CreateConnection
type CreateConnectionRequest struct {
	Bandwidth      string             `json:"bandwidth"`
	ConnectionName string             `json:"connectionName"`
	LagId          string             `json:"lagId"`
	Location       string             `json:"location"`
	ProviderName   string             `json:"providerName"`
	RequestMACSec  bool               `json:"requestMACSec"`
	Tags           []DirectConnectTag `json:"tags"`
}

// Connection is the response body for CreateConnection
// It is equivalent to Response of the action CreateConnection.
// https://docs.aws.amazon.com/directconnect/latest/APIReference/API_CreateConnection.html

type Connection struct {
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

// Implementing the Marshaler interface
func (c Connection) MarshalJSON() ([]byte, error) {
	type Alias Connection
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&c),
	})
}

// Implementing the Unmarshaler interface
func (c *Connection) UnmarshalJSON(b []byte) error {
	type Alias Connection
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	return json.Unmarshal(b, &aux)
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

type DeleteConnectionRequest struct {
	ConnectionId string `json:"connectionId"`
}

// DescribeConnectionsResponse is the response body for DescribeConnections
type DescribeConnectionsResponse struct {
	Connections []Connection `json:"connections"`
}

type UpdateConnectionRequest struct {
	ConnectionId   string `json:"connectionId"`
	ConnectionName string `json:"connectionName"`
	EncryptionMode string `json:"encryptionMode"`
}

func CreateConnection(r *http.Request) (Connection, error) {
	var dx CreateConnectionRequest
	var dc Connection

	// Unmarshal the body
	err := RequestToJson(r, &dx)
	if err != nil {
		return dc, err
	}

	// Assign the values
	dc.Bandwidth = dx.Bandwidth
	dc.ConnectionId = CreateConnectionID()
	dc.ConnectionName = dx.ConnectionName
	dc.ConnectionState = "available"
	dc.Location = dx.Location
	dc.ProviderName = dx.ProviderName
	dc.Tags = dx.Tags

	return dc, nil
}

func DescribeConnections(r *http.Request) (DescribeConnectionsRequest, error) {
	var request DescribeConnectionsRequest
	//var response DescribeConnectionsResponse
	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return request, err
	}
	return request, nil

	// response = DescribeConnectionsResponse{
	// 	Connections: []CreateConnectionResponse{dx},
	// }

}

func UpdateConnection(r *http.Request, dx *Connection) error {
	var request UpdateConnectionRequest

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return err
	}

	dx.ConnectionName = request.ConnectionName
	dx.EncryptionMode = request.EncryptionMode

	return nil
}

func DeleteConnection(r *http.Request, dx *Connection) error {
	var request DeleteConnectionRequest

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return err
	}

	dx.ConnectionState = "deleted"

	return nil
}

// CreateConnectionID generates a random connection ID
// The ID is prefixed with "dxcon-" followed by 8 random characters
func CreateConnectionID() string {
	return "dxcon-" + randomString(8)
}
