package pgsql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

//DB holds database connection pools
type DB struct {
	//Reader is readonly connection to database
	reader *pgxpool.Pool

	//Writer is read/write connection to database
	writer *pgxpool.Pool
}

//Conn provide reader or writer connection as per readonly state
func (db *DB) Conn(readonly bool) *pgxpool.Pool {
	if readonly {
		return db.reader
	}
	return db.writer
}

//DbConfig is config for disk persistent database (PostgreSQL)
type DbConfig struct {
	DbHost string
	DbPort uint16
	DbName string
	DbUser string
	DbPwd  string
	//DbTimeout is Connection timeout in seconds
	//if could not connect to server in given time then giveup and raise error
	DbTimeout int
	//DbSSLMode flag to enable disable SSL for database connection
	DbSSLMode string
}

//InitDbPool Initialize database connection ppol for PostgreSQL database
// Conection Options: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
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

//CloseDbPool close database connection. Not necessary to call, as pool itself closes inactive connections.
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

func initdb(config DbConfig, connName string) *pgxpool.Pool {
	var err error

	if config.DbPort < 1 {
		config.DbPort = 5432
	}
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=%s",
		config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName, config.DbTimeout, config.DbSSLMode)

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(fmt.Sprintf("pgsql connection parse error [%s]: %s", connName, err.Error()))
	}

	p, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		panic(fmt.Sprintf("pgsql connection failed [%s]: %s", connName, err.Error()))
	}
	return p
}
