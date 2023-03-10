package dx

import (
	"encoding/json"
	"math/rand"
	"net/http"
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
