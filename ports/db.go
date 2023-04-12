package ports

import (
	"encoding/json"
)

type DbPort interface {
	CloseDbConnection()
	GetVal(string, json.Unmarshaler) ([]byte, error)
	SetVal(string, json.Marshaler) error
}
