package databasepack

import (
	cryptpack "github.com/drempi/golang-REST-API/REST-API/cryptPack"
	errorpack "github.com/drempi/golang-REST-API/REST-API/errorPack"

	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

// DEFAULT ENTRIES IN USERS:
// login TEXT PRIMARY KEY: username
// password TEXT: a hashed password
// roles TEXT: a table of his roles

// Upon first run of server there are the following values:
// SUPERADMIN: PASSWORD, [0]
// VISITOR: PASSWORD, [2]

// Account that's what it is
type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Isnew    bool   `json:"isnew"`
}

// InitUsers initializes the users table
func InitUsers() {
	ExecCommand("CREATE TABLE IF NOT EXISTS users (login TEXT PRIMARY KEY, password TEXT, roles TEXT)")
	Password := cryptpack.CreateHash2("PASSWORD")
	// adding in SUPERADMIN
	ExecCommand("INSERT INTO users (login, password, roles) VALUES (\"SUPERADMIN\", \"" + Password + "\", \"[0]\")")
	// adding in VISITOR
	ExecCommand("INSERT INTO users (login, password, roles) VALUES (\"VISITOR\", \"" + Password + "\", \"[2]\")")
}

// FindUser checks if there is a user with given Login.
// If there is and the password is okay, return its roles
// If there is and the password is not okay, return -1
// If there is no such account, return -2
// If the list of roles is just plain wrong, return -3
func FindUser(acc Account) []int {
	command := "SELECT login, password, roles FROM users WHERE login = \"" + acc.Login + "\""
	row, err := D.Query(command)
	if errorpack.CHECK(&err) {
		return []int{-2}
	}
	defer row.Close()
	row.Next()
	var Login, Pass, Roles string
	err = row.Scan(&Login, &Pass, &Roles)
	if err != nil {
		return []int{-2}
	}
	err = bcrypt.CompareHashAndPassword([]byte(Pass), []byte(acc.Password))
	if err != nil {
		return []int{-1}
	}
	var roles []int
	err = json.Unmarshal([]byte(Roles), &roles)
	errorpack.OK(&err)
	return roles
}

// CreateUser creates an account
func CreateUser(acc Account) {
	acc.Password = cryptpack.CreateHash2(acc.Password)
	ExecCommand("INSERT INTO users (login, password, roles) VALUES (\"" + acc.Login + "\", \"" + acc.Password + "\", \"[2]\")")
}

// RemoveUser deletes an account
func RemoveUser(acc Account) {
	ExecCommand("DELETE FROM users WHERE login = \"" + acc.Login + "\"")
}
