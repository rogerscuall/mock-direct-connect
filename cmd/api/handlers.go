package main

import (
	"dx-mock/adapters/db"
	d "dx-mock/pkg/dx"
	"encoding/json"
	"log"
	"net/http"
)

func CreateConnection(w http.ResponseWriter, r *http.Request) {
	var err error

	dx, err = d.CreateConnection(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	connectionDB, err := db.NewAdapter(dbNameConnection)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()
	// Serialize the struct

	b, err := json.Marshal(dx)
	if err != nil {
		log.Println("Error serializing data", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	err = connectionDB.SetVal(dx.ConnectionId, b)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	// Update the tag database
	resourceTag := d.ResourceTag{
		ResourceArn: d.CreateARN("us-east-1", dx.ConnectionId),
		Tags:        dx.Tags,
	}

	tagDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer tagDB.CloseDbConnection()

	b, err = json.Marshal(resourceTag)
	if err != nil {
		log.Println("Error serializing data", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	err = tagDB.SetVal(dx.ConnectionId, b)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	// Return a response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dx)
}

func DescribeConnections(w http.ResponseWriter, r *http.Request) {
	request, err := d.DescribeConnections(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	connectionDB, err := db.NewAdapter(dbNameConnection)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	response := d.DescribeConnectionsResponse{
		Connections: []d.CreateConnectionResponse{},
	}

	// Find the connection in the database
	val, err := connectionDB.GetVal(request.ConnectionId)
	if err != nil {
		log.Println("Error in getting connection ID from database", err)
		json.NewEncoder(w).Encode(response)
	}

	// Unmarshal the data
	err = json.Unmarshal(val, &dx)
	if err != nil {
		log.Println("Error in unmarshalling data", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	response.Connections = append(response.Connections, dx)
	json.NewEncoder(w).Encode(response)
}

func DescribeTags(w http.ResponseWriter, r *http.Request) {
	request, err := d.DescribeTags(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	connectionDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	response := d.DescribeTagsResponse{}
	resourceTag := d.ResourceTag{}

	for _, resourceARN := range request.ResourceArns {
		key, err := d.GetIDFromARN(resourceARN)
		if err != nil {
			log.Println("Error in getting key from ARN", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		val, err := connectionDB.GetVal(key)
		if err != nil {
			log.Println("Value not found in database", err)
			continue
		}
		err = json.Unmarshal(val, &resourceTag)
		if err != nil {
			log.Println("Error in unmarshalling data", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		response.ResourceTags = append(response.ResourceTags, resourceTag)
	}

	json.NewEncoder(w).Encode(response)
}

func TagResource(w http.ResponseWriter, r *http.Request) {
	resourceTag, err := d.TagResource(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	connectionDB, err := db.NewAdapter(dbNameTags)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}
	defer connectionDB.CloseDbConnection()

	b, err := json.Marshal(resourceTag)
	if err != nil {
		log.Println("Error serializing data", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	key, err := d.GetIDFromARN(resourceTag.ResourceArn)
	if err != nil {
		log.Println("Error in getting key from ARN", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	err = connectionDB.SetVal(key, b)
	if err != nil {
		log.Println("Error in creating connection to database", err)
		http.Error(w, "Database Connection failure", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
