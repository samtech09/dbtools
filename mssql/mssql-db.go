package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/rs/zerolog"
)

//DB holds database connections
type DB struct {
	//Reader is readonly connection to database
	reader *sql.DB

	//Writer is read/write connection to database
	writer *sql.DB
}

//Conn provide reader or writer connection as per readonly state
func (db *DB) Conn(readonly bool) *sql.DB {
	if readonly {
		return db.reader
	}
	return db.writer
}

//DbConfig is config for disk persistent database
type DbConfig struct {
	DbHost string
	DbPort uint16
	DbName string
	DbUser string
	DbPwd  string
	//DbTimeout is Connection timeout in seconds
	//if could not connect to server in given time then giveup and raise error
	DbTimeout int
}

//InitDbPool Initialize database connection ppol for PostgreSQL database
func InitDbPool(reader, writer DbConfig) *DB {
	db := DB{}
	if reader.DbHost != "" {
		db.reader = initdb(reader, "sql-reader")
	}
	if writer.DbHost != "" {
		db.writer = initdb(writer, "sql-writer")
	}

	return &db
}

//CloseDbPool close database connection. Must be called to release connections before exiting from client.
func (db *DB) CloseDbPool() {
	if db.reader != nil {
		db.reader.Close()
	}
	if db.writer != nil {
		db.writer.Close()
	}
}

//SetLogger set zerolog logger for logging database events
func SetLogger(l zerolog.Logger) {
	// not implemented

	//cfg.LogLevel = pgx.LogLevelWarn
	//cfg.Logger = newLogger(l, connName)

	// pgxConnPoolConfig := pgxpool.Config{
	// 	ConnConfig:   cfg,
	// 	MaxConns:     8,
	// 	AfterConnect: nil,
	// }
}

func initdb(config DbConfig, connName string) *sql.DB {
	var err error
	var connString string
	var tout string
	var port string

	if config.DbTimeout > 0 {
		tout = fmt.Sprintf(";connection+timeout=%d", config.DbTimeout)
	}
	if config.DbPort > 0 {
		port = fmt.Sprintf(";port=%d", config.DbPort)
	}
	//log=1 log errors
	connString = fmt.Sprintf("server=%s%s;user id=%s;password=%s;database=%s%s;log=1", config.DbHost, port, config.DbUser, config.DbPwd, config.DbName, tout)
	//conn, err := sql.Open("mssql", connString)
	conn, err := sql.Open("sqlserver", connString)
	if err != nil {
		panic(fmt.Sprintf("mssql connection failed [%s]: %s", connName, err.Error()))
	}
	err = conn.Ping()
	if err != nil {
		panic(fmt.Sprintf("ping to mssql connection failed [%s]: %s", connName, err.Error()))
	}
	return conn
}
