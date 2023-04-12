package main

import (
	"dx-mock/adapters/db"
	"dx-mock/ports"
	"fmt"
)

func main() {
	var dbAdapter1 ports.DbPort
	dbAdapter1, err := db.NewAdapter("db1")
	if err != nil {
		panic(err)
	}
	defer dbAdapter1.CloseDbConnection()
	dbAdapter1.SetVal("key1", []byte("value1"))
	val, _ := dbAdapter1.GetVal("key1")
	fmt.Println(string(val))
}
