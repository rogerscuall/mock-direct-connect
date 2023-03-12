package dx

import (
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

func TagResource(r *http.Request) error {
	var request DirectConnectTag

	// Unmarshal the body
	err := RequestToJson(r, &request)
	if err != nil {
		return err
	}

	return nil
}

func CreateARN(region, id string) string {
	return fmt.Sprintf("arn:aws:directconnect:%s:123456789012:dxcon/%s", region, id)
}
