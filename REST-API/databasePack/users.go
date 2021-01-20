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
// email TEXT: email of user

// Upon first run of server there are the following values:
// SUPERADMIN: PASSWORD, [0], -
// VISITOR: PASSWORD, [2], -

// Account that's what it is
type Account struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	Additional string `json:"additional"`
	Isnew      bool   `json:"isnew"`
}

// InitUsers initializes the users table
func InitUsers() {
	ExecCommand("CREATE TABLE IF NOT EXISTS users (login TEXT PRIMARY KEY, password TEXT, roles TEXT, email TEXT UNIQUE)")
	Password := cryptpack.CreateHash2("PASSWORD")
	// adding in SUPERADMIN
	ExecCommand("INSERT INTO users (login, password, roles, email) VALUES (\"SUPERADMIN\", \"" + Password + "\", \"[0]\", \"empty1\")")
	// adding in VISITOR
	ExecCommand("INSERT INTO users (login, password, roles, email) VALUES (\"VISITOR\", \"" + Password + "\", \"[2]\", \"empty2\")")
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
	ExecCommand("INSERT INTO users (login, password, roles, email) VALUES (\"" + acc.Login + "\", \"" + acc.Password + "\", \"[2]\", \"" + acc.Additional + "\")")
}

// RemoveUser deletes an account
func RemoveUser(acc Account) {
	ExecCommand("DELETE FROM users WHERE login = \"" + acc.Login + "\"")
}

// ChangeGroup changes current groups of some user
func ChangeGroup(login string, tab []int) {
	newGroups, err := json.Marshal(tab)
	errorpack.OK(&err)
	ExecCommand("UPDATE users SET roles = \"" + string(newGroups) + "\" WHERE login = \"" + login + "\"")
}

// AddGroupUser adds group to user
func AddGroupUser(login string, group int) string {
	var acc Account
	acc.Login = login
	acc.Password = "aodpfksingrw"
	acc.Isnew = false
	tab := FindUser(acc)
	if tab[0] < -1 {
		return "User not found"
	}
	ADD := true
	for i := 0; i < len(tab); i++ {
		if tab[i] == group {
			ADD = false
		}
	}
	if !ADD {
		return "User already has this group"
	}
	tab = append(tab, group)
	ChangeGroup(login, tab)
	return "Groups changed successfully!"
}

// RemoveGroupUser removes group from user
func RemoveGroupUser(login string, group int) string {
	var acc Account
	acc.Login = login
	acc.Password = "aodpfksingrw"
	acc.Isnew = false
	tab := FindUser(acc)
	if tab[0] < -1 {
		return "User not found"
	}
	var newGroups []int
	REMOVE := false
	for i := 0; i < len(tab); i++ {
		if tab[i] == group {
			REMOVE = true
		} else {
			newGroups = append(newGroups, tab[i])
		}
	}
	if !REMOVE {
		return "User doesn't have this group"
	}
	ChangeGroup(login, newGroups)
	return "Groups changed successfully!"
}

// ChangePassword changes password of user
func ChangePassword(login string, newPassword string) {
	passwordCrypted := cryptpack.CreateHash2(newPassword)
	ExecCommand("UPDATE users SET password = \"" + passwordCrypted + "\" WHERE login = \"" + login + "\"")
}

// ResetPassword resets user's password to a new one
func ResetPassword(acc Account) string {
	roles := FindUser(acc)
	if roles[0] == -2 {
		return "No such username!"
	} else if roles[0] < 0 {
		return "Wrong Password!"
	}
	// Okay, thats how it looks now
	ChangePassword(acc.Login, acc.Additional)
	return "Password changed!"
}

// RandomPasswordUser changes the password of user to something random
func RandomPasswordUser(login string, hash string) string {
	command := "SELECT password FROM users WHERE login = \"" + login + "\""
	row, err := D.Query(command)
	if err != nil {
		return "No such username exists!"
	}
	defer row.Close()
	row.Next()
	var Pass string
	err = row.Scan(&Pass)
	errorpack.OK(&err)
	if Pass[0:len(Pass)/2] == hash {
		newPassword := cryptpack.RandomString(10)
		ChangePassword(login, newPassword)
		return "Password changed! new password: " + newPassword
	}
	return "Password not changed: unautharized query"
}

// SendEmail sends email to the person with this username
func SendEmail(login string) string {
	command := "SELECT email, password FROM users WHERE login = \"" + login + "\""
	row, err := D.Query(command)
	if err != nil {
		return "No such username exists!"
	}
	defer row.Close()
	row.Next()
	var Email, Pass string
	err = row.Scan(&Email, &Pass)
	errorpack.OK(&err)
	Pass = Pass[0 : len(Pass)/2]
	// Just the email!
	/*
		// Choose auth method and set it up
		auth := smtp.PlainAuth("", "piotr@mailtrap.io", "extremely_secret_pass", "smtp.mailtrap.io")

		// Here we do it all: connect to our server, set up a message and send it
		to := []string{"bill@gates.com"}
		msg := []byte("To: bill@gates.com\r\n" +
			"Subject: Why are you not using Mailtrap yet?\r\n" +
			"\r\n" +
			"Here’s the space for our great sales pitch\r\n")
		err = smtp.SendMail("smtp.mailtrap.io:25", auth, "piotr@mailtrap.io", to, msg)
		errorpack.OK(&err)
	*/
	return "Email sent!"
}
