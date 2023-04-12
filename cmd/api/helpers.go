package main

import (
	"encoding/json"
	"net/http"
)

func returnOk(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
