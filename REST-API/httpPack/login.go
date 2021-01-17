package httppack

import (
	cryptpack "github.com/REST-API/cryptPack"
	databasepack "github.com/REST-API/databasePack"

	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Account that what it is
type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Isnew    bool   `json:"isnew"`
}

// FindUser checks if there is a user with given Login.
// If there is and the password is okay, return its role
// If there is and the password is not okay, return -1
// If there is no such account, return -2
func FindUser(acc Account) int {
	command := "SELECT login, password, role FROM users WHERE login = \"" + acc.Login + "\""
	row, err := databasepack.D.Query(command)
	if err != nil {
		fmt.Println(err.Error())
		return -2
	}
	defer row.Close()
	row.Next()
	var Login, Pass string
	var Role int
	err = row.Scan(&Login, &Pass, &Role)
	if err != nil {
		return -2
	}
	err = bcrypt.CompareHashAndPassword([]byte(Pass), []byte(acc.Password))
	if err != nil {
		return -1
	}
	return Role
}

// CreateUser creates an account
func CreateUser(acc Account) {
	acc.Password = cryptpack.CreateHash2(acc.Password)
	statement, err := databasepack.D.Prepare("INSERT INTO users (login, password, role) VALUES (\"" + acc.Login + "\", \"" + acc.Password + "\", 1)")
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

// RemoveUser deletes an account
func RemoveUser(acc Account) {
	statement, err := databasepack.D.Prepare("DELETE FROM users WHERE login = \"" + acc.Login + "\"")
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}
