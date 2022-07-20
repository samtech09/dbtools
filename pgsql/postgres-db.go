package pgsql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

//DB holds database connection pools
type DB struct {
	//Reader is readonly connection to database
	reader DbConfig

	//Writer is read/write connection to database
	writer DbConfig
}

//Conn provide reader or writer connection as per readonly state
func (db *DB) Conn(readonly bool) *pgx.Conn {
	if readonly {
		return getDbConn(db.reader, "sql-reader")
	}
	return getDbConn(db.writer, "sql-writer")
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
	DbSSLMode  string
	DbPoolSize uint16
}

//InitDb Initialize database connection and mke sure database is reachable
// Conection Options: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func InitDb(reader, writer DbConfig) *DB {
	db := DB{}
	if reader.DbHost != "" {
		c := getDbConn(reader, "sql-reader")
		defer c.Close(context.Background())
		err := c.Ping(context.Background())
		if err != nil {
			panic(fmt.Sprintf("pgsql connection PING failed [%s]: %s", "sql-reader", err.Error()))
		}
	}
	if writer.DbHost != "" {
		c := getDbConn(writer, "sql-writer")
		defer c.Close(context.Background())
		err := c.Ping(context.Background())
		if err != nil {
			panic(fmt.Sprintf("pgsql connection PING failed [%s]: %s", "sql-writer", err.Error()))
		}
	}
	db.reader = reader
	db.writer = writer
	return &db
}

func getDbConn(config DbConfig, connName string) *pgx.Conn {
	var err error

	if config.DbPort < 1 {
		config.DbPort = 5432
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=%s&statement_cache_mode=describe",
		config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName, config.DbTimeout, config.DbSSLMode)

	cfg, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(fmt.Sprintf("pgsql connection parse error [%s]: %s", connName, err.Error()))
	}

	if config.DbPoolSize > 0 {
		cfg.PreferSimpleProtocol = true
	}

	c, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		panic(fmt.Sprintf("pgsql connection failed [%s]: %s", connName, err.Error()))
	}
	return c
}

func BulkCopy(cntxt context.Context, conn *pgx.Conn, targetTable string, columns []string, rows [][]interface{}) (int64, error) {
	return conn.CopyFrom(cntxt,
		pgx.Identifier{targetTable},
		columns,
		pgx.CopyFromRows(rows),
	)
}
