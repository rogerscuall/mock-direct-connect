package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var dx *CreateConnectionResponse

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Mock Direct Connect API server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

type DescribeConnectionsRequest struct {
	ConnectionId string `json:"connectionId"`
}

// CreateConnectionRequest
type CreateConnectionRequest struct {
	Bandwidth      string `json:"bandwidth"`
	ConnectionName string `json:"connectionName"`
	LagId          string `json:"lagId"`
	Location       string `json:"location"`
	ProviderName   string `json:"providerName"`
	RequestMACSec  bool   `json:"requestMACSec"`
	Action         string `json:"Action"`
}

// CreateConnectionResponse
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

type DirectConnectKey struct {
	Ckn       string `json:"ckn"`
	SecretArn string `json:"secretARN"`
	StartOn   string `json:"startOn"`
	State     string `json:"state"`
}

type DirectConnectTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DescribeConnectionsResponse struct {
	Connections []CreateConnectionResponse `json:"connections"`
}

type DescribeTagsRequest struct {
	ResourceArns []string `json:"resourceArns"`
}

type ResourceTag struct {
	ResourceArn string `json:"resourceArn"`
	Tags        []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
}

type DescribeTagsResponse struct {
	ResourceTags []ResourceTag `json:"resourceTags"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Get the Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-amz-json-1.1" {
		log.Println("Content-Type is not application/x-amz-json-1.1")
	}
	// Get the target
	target := r.Header.Get("X-Amz-Target")
	log.Println("X-Amz-Target is ", target)
	switch target {
	case "OvertureService.CreateConnection":
		var err error

		dx, err = CreateConnection(r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Return a response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dx)

		return
	case "OvertureService.DescribeConnections":
		response, err := DescribeConnections(r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(response)

		return
	case "OvertureService.DescribeTags":
		response, err := DescribeTags(r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(response)

		return
	}

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

func DescribeConnections(r *http.Request) (*DescribeConnectionsResponse, error) {
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

func DescribeTags(r *http.Request) (*DescribeTagsResponse, error) {
	var request DescribeTagsRequest

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return nil, err
	}

	response := DescribeTagsResponse{
		ResourceTags: []ResourceTag{
			{
				ResourceArn: "arn:aws:directconnect:us-east-1:1234567890:dxcon-1234567890",
			},
		},
	}
	return &response, nil
}
