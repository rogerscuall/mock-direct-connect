package db

import (
	"encoding/json"
	"log"

	"github.com/charmbracelet/charm/kv"
)

type Adapter struct {
	db *kv.KV
}

// NewAdapter returns a new Adapter
// It will connect to the DB and Sync to the latest update
func NewAdapter(dbName string) (*Adapter, error) {
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		return nil, err
	}
	db.Sync()
	return &Adapter{db: db}, nil
}

func (da Adapter) CloseDbConnection() {
	err := da.db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}

// GetVal will update the struct s with the value from the DB
// It will return an error.
func (da Adapter) GetVal(key string, s json.Unmarshaler) error {
	val, err := da.db.Get([]byte(key))
	if err != nil {
		return err
	}
	err = s.UnmarshalJSON(val)
	if err != nil {
		return err
	}
	return nil
}

// SetVal will set the value of the key in the DB
// It will return an error.
func (da Adapter) SetVal(key string, s json.Marshaler) error {
	val, err := s.MarshalJSON()
	if err != nil {
		return err
	}
	err = da.db.Set([]byte(key), val)
	if err != nil {
		return err
	}
	return nil
}

func (da Adapter) Sync() {
	err := da.db.Sync()
	if err != nil {
		log.Fatalf("Error syncing database: %v", err)
	}
}
