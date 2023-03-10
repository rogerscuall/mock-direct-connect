package dx

import "net/http"

// DescribeTagsRequest is the request body for DescribeTags
type DescribeTagsRequest struct {
	ResourceArns []string `json:"resourceArns"`
}

type DescribeTagsResponse struct {
	ResourceTags []ResourceTag `json:"resourceTags"`
}

type ResourceTag struct {
	ResourceArn string `json:"resourceArn"`
	Tags        []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
}

type DirectConnectTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
