package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Open(host string, port int, dbname, user, pwd string, ssl bool) (*sql.DB, error) {
	sslEnabled := func() string {
		if ssl {
			return "enable"
		}
		return "disable"
	}
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, pwd, dbname, sslEnabled())
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err // consider wrapping it providing more information about the origins ofthe error
	}
	// checking if the database is actually available
	// we can do it by means of the db.Ping() function or actually querying the DB, ;a simple
	// query, something like select 1 will suffice
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
