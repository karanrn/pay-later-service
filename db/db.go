package db

import (
	"database/sql"
	"os"
)

// Connect returns connection item for the database
func Connect() (db *sql.DB) {
	// DB parameters
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	dbName := os.Getenv("DBNAME")
	port := os.Getenv("DBPORT")
	host := os.Getenv("DBHOST")
	dbDriver := "mysql"

	db, err := sql.Open(dbDriver, user+":"+password+"@tcp("+host+":"+port+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}
