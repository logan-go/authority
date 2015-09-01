package authority

import (
	"database/sql"
	"errors"
)

var authority_db *sql.DB
var authority_prefix string
var authority_database_name string

func SetDB(db *sql.DB) error {
	authority_db = db
	return CheckDBConn()
}

func SetPrefix(prefix string) {
	authority_prefix = prefix
}

func GetPrefix() string {
	return authority_prefix
}

func SetDatabaseName(databaseName string) {
	authority_database_name = databaseName
}

func GetDatabaseName() string {
	return authority_database_name
}

func CheckDBConn() error {
	if authority_db == nil {
		return errors.New("Connection is Empty.")
	}
	return authority_db.Ping()
}
