package databasepack

import (
	"database/sql"

	errorpack "github.com/drempi/golang-REST-API/REST-API/errorPack"
)

// DEFAULT ENTRIES IN DISPATCHER:
// name TEXT PRIMARY KEY: name of given table
// created_at DATETIME: time and date when this table was created
// updated_at DATETIME: time and date when this table was updated

// Upon first run of server there are the following values:
// dispatcher: ??, ??
// licenser: ??, ??
// users: ??, ??

// TableType its a format of added table
type TableType struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// InitDispatcher it initializes the dispatcher table
func InitDispatcher() {
	ExecCommand("CREATE TABLE IF NOT EXISTS dispatcher (name TEXT PRIMARY KEY, type INTEGER, created_at DATETIME, updated_at DATETIME)")
	// adding dispatcher
	ExecCommand("INSERT INTO dispatcher (name, type, created_at, updated_at) VALUES (\"dispatcher\", 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	// adding licenser
	ExecCommand("INSERT INTO dispatcher (name, type, created_at, updated_at) VALUES (\"licenser\", 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	// adding users
	ExecCommand("INSERT INTO dispatcher (name, type, created_at, updated_at) VALUES (\"users\", 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
}

// FindTable checks if table with such name exists.
func FindTable(name string) int {
	var T int
	command := "SELECT type FROM dispatcher WHERE name = \"" + name + "\""
	err := D.QueryRow(command).Scan(&T)
	if err != nil {
		if err != sql.ErrNoRows {
			errorpack.OK(&err)
		}
		return -1
	}
	return T
}

// CheckName checks if []byte has a proper name for a table
// It must have at least one lower case character (so that sql doesn't break)
// It must be made only out of english alphabet
func CheckName(name []byte) bool {
	for i := 0; i < len(name); i++ {
		if name[i] <= 64 || name[i] > 122 {
			return false
		} else if name[i] > 90 && name[i] <= 96 {
			return false
		}
	}
	for i := 0; i < len(name); i++ {
		if name[i] > 96 {
			return true
		}
	}
	return false
}

// AddTable adds empty table to the database
func AddTable(tab TableType) {
	if tab.Type == 1 {
		// This is the simpler type
		ExecCommand("INSERT INTO dispatcher (name, type, created_at, updated_at) VALUES (\"" + tab.Name + "\", 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
		ExecCommand("CREATE TABLE " + tab.Name + " (id INTEGER PRIMARY KEY AUTOINCREMENT, data TEXT, created_at DATETIME, updated_at DATETIME)")
		//ExecCommand("INSERT INTO " + tab.Name + " (id, data, created_at, updated_at) VALUES (0, \"nothing\", CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
		AddDefaultPermissions(tab.Name)
	} else if tab.Type == 2 {
		// This is the double type
		ExecCommand("INSERT INTO dispatcher (name, type, created_at, updated_at) VALUES (\"" + tab.Name + "\", 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
		ExecCommand("CREATE TABLE " + tab.Name + " (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME, updated_at DATETIME)")
		//ExecCommand("INSERT INTO " + tab.Name + " (id, created_at, updated_at) VALUES (0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
		AddDefaultPermissions(tab.Name)
		ExecCommand("INSERT INTO dispatcher (name, type, created_at, updated_at) VALUES (\"" + tab.Name + "_lang\", 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
		ExecCommand("CREATE TABLE " + tab.Name + "_lang (id INTEGER NOT NULL, lang INTEGER NOT NULL, data TEXT, PRIMARY KEY(id, lang))")
		AddDefaultPermissions(tab.Name + "_lang")
	}
}
