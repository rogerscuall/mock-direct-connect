package dx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DescribeTagsRequest is the request body for DescribeTags
type DescribeTagsRequest struct {
	ResourceArns []string `json:"resourceArns"`
}

type DescribeTagsResponse struct {
	ResourceTags []ResourceTag `json:"resourceTags"`
}

type ResourceTag struct {
	ResourceArn string             `json:"resourceArn"`
	Tags        []DirectConnectTag `json:"tags"`
}

// Implement the Marshaler interface
func (r ResourceTag) MarshalJSON() ([]byte, error) {
	var Alias struct {
		ResourceArn string             `json:"resourceArn"`
		Tags        []DirectConnectTag `json:"tags"`
	}
	Alias.ResourceArn = r.ResourceArn
	Alias.Tags = r.Tags
	return json.Marshal(&Alias)
}

// Implement the Unmarshaler interface
func (r *ResourceTag) UnmarshalJSON(b []byte) error {
	var Alias struct {
		ResourceArn string             `json:"resourceArn"`
		Tags        []DirectConnectTag `json:"tags"`
	}
	err := json.Unmarshal(b, &Alias)
	if err != nil {
		return err
	}
	r.ResourceArn = Alias.ResourceArn
	r.Tags = Alias.Tags
	return nil
}

type DirectConnectTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func DescribeTags(r *http.Request) (DescribeTagsRequest, error) {
	var request DescribeTagsRequest

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return request, err
	}

	if len(request.ResourceArns) == 0 {
		return request, fmt.Errorf("ResourceArns is empty")
	}

	return request, nil
}

func TagResource(r *http.Request) (ResourceTag, error) {
	var request ResourceTag

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return request, err
	}

	return request, nil
}

func CreateARN(region, id string) string {
	return fmt.Sprintf("arn:aws:directconnect:%s:123456789012:dxcon/%s", region, id)
}
