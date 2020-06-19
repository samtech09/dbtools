package pgsql

import (
	"context"
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	// port 0 will tell driver to use default port
	cfg := DbConfig{"192.168.60.206", 0, "testdb", "testuser", "testuser", 30, "disable", 10}
	db := InitDb(cfg, DbConfig{})
	c := db.Conn(true)
	defer c.Close(context.Background())
	rows, err := c.Query(context.Background(), "select id, email from users limit 1;")

	if err != nil {
		t.Error(err)
		return
	}

	var id int
	var title string
	if rows.Next() {
		err := rows.Scan(&id, &title)
		if err != nil {
			t.Error(err)
			return
		}
	}

	fmt.Printf("ID: %d, Title: %s\n", id, title)

	if id <= 0 {
		t.Error("Id must be > 0")
		return
	}
}
