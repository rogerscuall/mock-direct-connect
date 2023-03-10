package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	d "dx-mock/pkg/dx"
)

var dx d.CreateConnectionResponse

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

		dx, err = d.CreateConnection(r)
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
		response, err := d.DescribeConnections(r, dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(response)

		return
	case "OvertureService.DescribeTags":
		response, err := d.DescribeTags(r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(response)
		return
	case "OvertureService.DeleteConnection":
		err := d.DeleteConnection(r, &dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(dx)

		return
	case "OvertureService.UpdateConnection":
		err := d.UpdateConnection(r, &dx)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(dx)

		return
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Mock Direct Connect API server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
