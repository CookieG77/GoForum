package functions

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var databaseInitialised = false
var db *sql.DB

// InitDatabase initialises the database connection
func InitDatabase() {
	if !databaseInitialised {
		if os.Getenv("DB_URL") == "" {
			ErrorPrintf("DB_URL environment variable not set\n")
			return
		}
		testDB, err := sql.Open(os.Getenv("DB_URL"), os.Getenv("DB_NAME"))
		if err != nil {
			ErrorPrintf("Error opening database: %v\n", err)
			return
		}
		err = testDB.Ping()
		if err != nil {
			ErrorPrintf("Error pinging database: %v\n", err)
			return
		}
		db = testDB
		InfoPrintf("Database initialised\n")
		databaseInitialised = true
	}
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if databaseInitialised {
		InfoPrintf("Database closed\n")
		err := db.Close()
		if err != nil {
			ErrorPrintf("Error closing database: %v\n", err)
			return
		}
		databaseInitialised = false
	}
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// IsDatabaseInitialised returns whether the database connection is initialised
func IsDatabaseInitialised() bool {
	return databaseInitialised
}
