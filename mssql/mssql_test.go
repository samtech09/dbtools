package mssql

import (
	"database/sql"
	"fmt"
	"testing"
)

var db *DB

func TestConnection(t *testing.T) {
	cfg := DbConfig{"172.25.12.203", 0, "testdb", "testuser", "testuser", 30, "disable"}
	db = InitDbPool(cfg, cfg)

	fmt.Println("no panic! connection suceeded")
}

func TestOps(t *testing.T) {
	sql := "IF OBJECT_ID('dbo.Students') IS NOT NULL drop TABLE dbo.Students;"
	_, err := db.Conn(false).Exec(sql)
	if err != nil {
		t.Errorf("Drop table error: %s", err.Error())
		t.FailNow()
	}

	sql = `create table Students (
		id int not null primary key,
		Name varchar(32) not null,
		Age int not null);`
	_, err = db.Conn(false).Exec(sql)
	if err != nil {
		t.Errorf("Create table error: %s", err.Error())
		t.FailNow()
	}

	// insert
	sql = "Insert into Students(id, name, age) values(@p1, @p2, @p3);"
	_, err = db.Conn(false).Exec(sql, 1, "User1", 12)
	if err != nil {
		t.Errorf("Insert error: %s", err.Error())
		t.FailNow()
	}
	_, err = db.Conn(false).Exec(sql, 2, "User2", 9)
	if err != nil {
		t.Errorf("Insert error: %s", err.Error())
		t.FailNow()
	}

	// select
	sql = "select * from students where age>@p1"
	rows, err := db.Conn(true).Query(sql, 10)
	if err != nil {
		t.Errorf("Select error: %s", err.Error())
		t.FailNow()
	}
	// get count of rows
	cnt := countRows(rows)
	exp := 1
	if cnt != exp {
		t.Errorf("Select error. Expected: %d,  Got: %d", exp, cnt)
		t.FailNow()
	}
}

func TestProc(t *testing.T) {
	//drop proc if exist
	sql := "IF OBJECT_ID('dbo.filterstud') IS NOT NULL drop PROC dbo.filterstud;"
	_, err := db.Conn(false).Exec(sql)
	if err != nil {
		t.Errorf("Drop table error: %s", err.Error())
		t.FailNow()
	}
	//create proc
	sql = `CREATE PROCEDURE filterstud
			@age int
		AS BEGIN
			select * from students where age>@age
		END;`
	_, err = db.Conn(false).Exec(sql)
	if err != nil {
		t.Errorf("Create proc error: %s", err.Error())
		t.FailNow()
	}

	//execute proc
	// select
	sql = "exec filterstud @p1"
	rows2, err := db.Conn(true).Query(sql, 8)
	if err != nil {
		t.Errorf("Exec proc error: %s", err.Error())
		t.FailNow()
	}
	// get count of rows
	cnt := countRows(rows2)
	exp := 2
	if cnt != exp {
		t.Errorf("Expec proc error. Expected: %d,  Got: %d", exp, cnt)
		t.FailNow()
	}
}

func TestClose(t *testing.T) {
	db.CloseDbPool()
}

func countRows(rw *sql.Rows) int {
	cnt := 0
	defer rw.Close()
	for rw.Next() {
		cnt++
	}
	return cnt
}
