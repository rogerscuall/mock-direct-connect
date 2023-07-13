package dx

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// randomString generates a random string of length n
func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func RequestToJson(r *http.Request, v interface{}) error {
	// Unmarshal the body
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}

	return nil
}

func GetIDFromARN(arn string) (string, error) {
	splitArn := strings.Split(arn, "/")
	if len(splitArn) != 2 {
		return "", fmt.Errorf("invalid ARN: %s", arn)
	}
	return splitArn[1], nil
}

// createDxGatewayAssociationID create a random string representing the DXGatewayAssociation ID
// Example is : 86bb6da8-c587-4b3d-89a5-f8335defc5ad
func createDxGatewayAssociationID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		randomString(8),
		randomString(4),
		randomString(4),
		randomString(4),
		randomString(12),
	)
}