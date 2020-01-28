# dbtools
Database modules to connect PostgreSQL, SQL-Server or MongoDB.

Each tools is self-contained module.

## Features
- Easy to use plug and play type modules.
- Provide query splitting (read/write) for Postgres and SQL-Server.

## Usage
Import, provide config and start using.

```
import "github.com/samtech09/dbtools/mssql"
...

func InitConnection() {
	cfg := DbConfig{"127.0.0.1", 0, "testdb", "testuser", "testuser", 30}
	db = InitDbPool(cfg, cfg)
	fmt.Println("Connection suceeded")
}

func TestOps(t *testing.T) {
	sql = "select * from table1 where age>@p1"
	rows, err := db.Conn(true).Query(sql, 10)
	if err != nil {
		t.Errorf("Select error: %s", err.Error())
		t.FailNow()
	}
	...
  // scan rows...
  ...
}

```

## Technologies used
- [pgx](https://github.com/jackc/pgx) and [pgxpool](https://github.com/jackc/pgx/v4/pgxpool) for connection to PostgreSQL
- [go-mssqldb](https://github.com/denisenkom/go-mssqldb) for connection to MsSQL Server
- [official mongodb driver](https://go.mongodb.org/mongo-driver/mongo) for working with MongoDB
- [zerolog](https://github.com/rs/zerolog) for logging
