package databasepack

import (
	errorpack "github.com/drempi/golang-REST-API/REST-API/errorPack"

	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type request struct {
	ID        int    `json:"id"`
	Lang      int    `json:"lang"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"udpated_at"`
	Data      string `json:"data"`
}

type multirequest struct {
	Requests []request `json:"requests"`
}

// D its a referenve to database
var D *sql.DB

// Init initializes database
func Init() {
	DATABASE, err := os.Open("DB.db")
	newdatabase := false
	if err != nil {
		newdatabase = true
		DATABASE, err = os.Create("DB.db")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	DATABASE.Close()
	log.Println("Database ready.")

	D, _ = sql.Open("sqlite3", "DB.db")

	if newdatabase {
		InitDispatcher()
		InitLicenser()
		InitUsers()
	}
}

// ExecCommand just a shorthand for executing a command
func ExecCommand(s string) {
	statement, err := D.Prepare(s)
	errorpack.OK(&err)
	_, err = statement.Exec()
	errorpack.OK(&err)
}

// UpdateDate update the updated_at field
func UpdateDate(name string, field string, id string) {
	ExecCommand("UPDATE " + name + " SET updated_at = CURRENT_TIMESTAMP WHERE " + field + " = " + id)
}

// CreatePosts creates posts in given table
func CreatePosts(text []byte, name string, T int) {
	var R multirequest
	err := json.Unmarshal(text, &R)
	if errorpack.CHECK(&err) {
		return
	}

	if T == 1 {
		for i := 0; i < len(R.Requests); i++ {
			if R.Requests[i].ID < 0 {
				ExecCommand("INSERT INTO " + name + " (data, created_at, updated_at) VALUES (\"" + R.Requests[i].Data + "\", CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
			} else {
				ExecCommand("UPDATE " + name + " SET data = " + R.Requests[i].Data + ", updated_at = CURRENT_TIMESTAMP WHERE id = " + strconv.Itoa(R.Requests[i].ID))
			}
		}
	} else if T == 2 {
		for i := 0; i < len(R.Requests); i++ {
			if R.Requests[i].ID < 0 {
				ExecCommand("INSERT INTO " + name + " (created_at, updated_at) VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
			} else {
				ExecCommand("UPDATE " + name + " SET updated_at = CURRENT_TIMESTAMP WHERE id = " + strconv.Itoa(R.Requests[i].ID))
				if ExistsPost(name, R.Requests[i].ID, R.Requests[i].Lang) {
					// if this row already exists
					ExecCommand("UPDATE " + name + "_lang SET data = " + R.Requests[i].Data + " WHERE id = " + strconv.Itoa(R.Requests[i].ID) + " AND lang = " + strconv.Itoa(R.Requests[i].Lang))
				} else {
					// if this row doesn't exist yet
					ExecCommand("INSERT INTO " + name + "_lang (id, lang, data) VALUES (" + strconv.Itoa(R.Requests[i].ID) + ", " + strconv.Itoa(R.Requests[i].Lang) + ", " + R.Requests[i].Data + ")")
				}
			}
		}
	}
	UpdateDate("dispatcher", "name", "\""+name+"\"")
}

// ExistsPost checks if post exists
func ExistsPost(name string, id int, lang int) bool {
	var I int
	var command string
	if lang < 0 {
		command = "SELECT id FROM " + name + " WHERE id = " + strconv.Itoa(id)
	} else {
		command = "SELECT id FROM " + name + "_lang WHERE id = " + strconv.Itoa(id) + " AND lang = " + strconv.Itoa(lang)
	}
	err := D.QueryRow(command).Scan(&I)
	if err != nil {
		if err != sql.ErrNoRows {
			errorpack.OK(&err)
		}
		return false
	}
	return true
}

// GetPosts gives the ordered posts
func GetPosts(name string, T int, offset int, amt int, lang int) []byte {
	var R multirequest
	if T == 1 {
		command := "SELECT * FROM " + name + " LIMIT " + strconv.Itoa(amt) + " OFFSET " + strconv.Itoa(offset)
		rows, err := D.Query(command)
		if errorpack.CHECK(&err) {
			return []byte("Failed to get rows")
		}
		defer rows.Close()
		var newRow request
		var createdAt, updatedAt, data string
		var id int
		newRow.Lang = -1
		for rows.Next() {
			err = rows.Scan(&id, &data, &createdAt, &updatedAt)
			if err != nil {
				break
			}
			newRow.ID = id
			newRow.Data = data
			newRow.CreatedAt = createdAt
			newRow.UpdatedAt = updatedAt
			R.Requests = append(R.Requests, newRow)
		}
	} else if T == 2 {
		command := "SELECT * FROM " + name + " LIMIT " + strconv.Itoa(amt) + " OFFSET " + strconv.Itoa(offset)
		rows, err := D.Query(command)
		if errorpack.CHECK(&err) {
			return []byte("Failed to get rows")
		}
		defer rows.Close()
		var newRow request
		var createdAt, updatedAt, data string
		var id, lang int
		for rows.Next() {
			err = rows.Scan(&id, &createdAt, &updatedAt)
			if err != nil {
				break
			}
			newRow.ID = id
			newRow.CreatedAt = createdAt
			newRow.UpdatedAt = updatedAt
			command2 := "SELECT * FROM " + name + "_lang WHERE id = " + strconv.Itoa(id)
			err2 := D.QueryRow(command2).Scan(&id, &lang, &data)
			if err2 == nil {
				newRow.Lang = lang
				newRow.Data = data
				R.Requests = append(R.Requests, newRow)
			} else if err2 != sql.ErrNoRows {
				errorpack.OK(&err)
			}
		}
	}
	ANS, err := json.Marshal(R)
	errorpack.OK(&err)
	return ANS
}
