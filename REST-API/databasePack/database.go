package databasepack

import (
	"database/sql"
	"log"
	"os"
)

// D its a referenve to database
var D *sql.DB

// Init initializes database
func Init() {
	DATABASE, err := os.Open("table.db")
	if err != nil {
		DATABASE, err = os.Create("table.db")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	DATABASE.Close()
	log.Println("Database created.")

	D, _ = sql.Open("sqlite3", "table.db")

	statement, _ := D.Prepare("CREATE TABLE IF NOT EXISTS users (login TEXT PRIMARY KEY, password TEXT, role INTEGER)")
	statement.Exec()
}

// There will be more in the future
