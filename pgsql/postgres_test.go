package pgsql

import (
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	// port 0 will tell driver to use default port
	cfg := DbConfig{"192.168.60.206", 0, "testdb", "testuser", "testuser", 30, "disable"}
	db := InitDbPool(cfg, DbConfig{})
	defer db.CloseDbPool()

	fmt.Println("no panic! connection suceeded")
}
