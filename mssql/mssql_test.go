package mssql

import (
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	cfg := DbConfig{"172.25.12.203", 0, "testdb", "testuser", "testuser", 30}
	db := InitDbPool(cfg, DbConfig{})
	defer db.CloseDbPool()

	fmt.Println("no panic! connection suceeded")
}
