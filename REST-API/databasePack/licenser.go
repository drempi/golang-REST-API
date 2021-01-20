package databasepack

import (
	errorpack "github.com/drempi/golang-REST-API/REST-API/errorPack"

	"strconv"
)

// DEFAULT ENTRIES IN LICENSER:
// group INT PRIMARY KEY: index of group
// creative BOOL: can he create/remove tables?
// viewer_table BOOL: can he view table?
// editor_table BOOL: can he edit table?
// licenser_table BOOL: can he change the permissions for table?

// Upon first run of server there are the following entries:
// group,
// creator,
// viewer_dispatcher, editor_dispatcher, licenser_dispatcher,
// viewer_licenser, editor_licenser, licenser_licenser,
// viewer_users, editor_users, licenser_users

// And the following values:

// 0 (superadmin): 1, 1, 1, 1, 1, 1, 1, 1, 1, 1
// 1 (admin): 1, 1, 0, 1, 0, 1, 0, 1, 0
// 2 (visitor): 0, 0, 0, 0, 0, 0, 0, 0, 0

// InitLicenser it initializes the licenser table
func InitLicenser() {
	ExecCommand("CREATE TABLE IF NOT EXISTS licenser (id INTEGER PRIMARY KEY, creator_ BOOL DEFAULT 0, licenser_ BOOL DEFAULT 0, viewer_dispatcher BOOL DEFAULT 0, editor_dispatcher BOOL DEFAULT 0, viewer_licenser BOOL DEFAULT 0, editor_licenser BOOL DEFAULT 0, viewer_users BOOL DEFAULT 0, editor_users BOOL DEFAULT 0)")
	// adding 0 group (superadmin)
	ExecCommand("INSERT INTO licenser (id, creator_, licenser_, viewer_dispatcher, editor_dispatcher, viewer_licenser, editor_licenser, viewer_users, editor_users) VALUES (0, 1, 1, 1, 1, 1, 1, 1, 1)")
	// adding 1 group (visitor)
	ExecCommand("INSERT INTO licenser (id, creator_, licenser_, viewer_dispatcher, editor_dispatcher, viewer_licenser, editor_licenser, viewer_users, editor_users) VALUES (1, 1, 0, 1, 0, 1, 0, 1, 0)")
	// adding 2 group (admin)
	ExecCommand("INSERT INTO licenser (id, creator_, licenser_, viewer_dispatcher, editor_dispatcher, viewer_licenser, editor_licenser, viewer_users, editor_users) VALUES (2, 0, 0, 0, 0, 0, 0, 0, 0)")
}

// Allowed checks if the user is able to perform given action
func Allowed(roles []int, action string) bool {
	var converted string
	converted = "("
	for i := 0; i < len(roles); i++ {
		converted = converted + strconv.Itoa(roles[i])
		if i == len(roles)-1 {
			converted = converted + ")"
		} else {
			converted = converted + ", "
		}
	}
	command := "SELECT " + action + " FROM licenser WHERE id IN " + converted
	row, err := D.Query(command)
	if errorpack.CHECK(&err) {
		return false
	}
	defer row.Close()
	var Able bool
	for row.Next() {
		err = row.Scan(&Able)
		if err != nil {
			return false
		}
		if Able {
			return true
		}
	}
	return false
}

// AddPermission Adds permission for given action
// on default its just 1 for superadmin and 0 for everyone else
func AddPermission(action string) {
	ExecCommand("ALTER TABLE licenser ADD " + action + " BOOL DEFAULT 0")
	ExecCommand("UPDATE licenser SET " + action + " = 1 WHERE id = 0")
}

// AddDefaultPermissions adds 3 default permissions
func AddDefaultPermissions(name string) {
	AddPermission("viewer_" + name)
	AddPermission("editor_" + name)
}

// EditLicenser changes the licenser table
func EditLicenser(group int, action string, val bool) {
	if val {
		ExecCommand("UPDATE licenser SET " + action + " = 1 WHERE id = " + strconv.Itoa(group))
	} else {
		ExecCommand("UPDATE licenser SET " + action + " = 0 WHERE id = " + strconv.Itoa(group))
	}
}

// AddGroup adds a new empty group with certain id and with no licenses
func AddGroup(id int) string {
	command := "SELECT id FROM licenser WHERE id = " + strconv.Itoa(id)
	_, err := D.Query(command)
	if err == nil {
		return "id already taken"
	} else if id < 0 {
		return "id has to be positive"
	}
	ExecCommand("INSERT licenser (id) VALUES (" + strconv.Itoa(id) + ")")
	return "completed!"
}
