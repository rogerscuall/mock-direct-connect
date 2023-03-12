package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	d "dx-mock/pkg/dx"
)

var (
	dx               d.CreateConnectionResponse
	dbNameConnection = "connection"
	dbNameTags       = "tags"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Get the Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-amz-json-1.1" {
		log.Println("Content-Type is not application/x-amz-json-1.1")
	}
	// Get the target
	serviceAction := strings.Split(r.Header.Get("X-Amz-Target"), ".")
	if len(serviceAction) != 2 {
		log.Println("X-Amz-Target is not in the correct format")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	service := serviceAction[0]
	if service != "OvertureService" {
		log.Println("Service is not OvertureService")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	action := serviceAction[1]
	log.Println("Request for: ", action)
	switch action {
	case "CreateConnection":
		CreateConnection(w, r)
	case "DescribeConnections":
		DescribeConnections(w, r)
	case "DescribeTags":
		DescribeTags(w, r)
	case "DeleteConnection":
		err := d.DeleteConnection(r, &dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(dx)

		return
	case "UpdateConnection":
		err := d.UpdateConnection(r, &dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(dx)

		return
	case "TagResource":
		TagResource(w, r)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Mock Direct Connect API server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
