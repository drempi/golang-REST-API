package databasepack

import (
	"strconv"

	errorpack "github.com/drempi/golang-REST-API/REST-API/errorPack"
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
// 1 (admin): 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
// 2 (visitor): 1, 1, 0, 0, 1, 0, 0, 1, 0, 0

// InitLicenser it initializes the licenser table
func InitLicenser() {
	ExecCommand("CREATE TABLE IF NOT EXISTS licenser (id INTEGER PRIMARY KEY, creator_ BOOL, viewer_dispatcher BOOL, editor_dispatcher BOOL, licenser_dispatcher BOOL, viewer_licenser BOOL, editor_licenser BOOL, licenser_licenser BOOL, viewer_users BOOL, editor_users BOOL, licenser_users BOOL)")
	// adding 0 group (superadmin)
	ExecCommand("INSERT INTO licenser (id, creator_, viewer_dispatcher, editor_dispatcher, licenser_dispatcher, viewer_licenser, editor_licenser, licenser_licenser, viewer_users, editor_users, licenser_users) VALUES (0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)")
	// adding 1 group (visitor)
	ExecCommand("INSERT INTO licenser (id, creator_, viewer_dispatcher, editor_dispatcher, licenser_dispatcher, viewer_licenser, editor_licenser, licenser_licenser, viewer_users, editor_users, licenser_users) VALUES (1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0)")
	// adding 2 group (admin)
	ExecCommand("INSERT INTO licenser (id, creator_, viewer_dispatcher, editor_dispatcher, licenser_dispatcher, viewer_licenser, editor_licenser, licenser_licenser, viewer_users, editor_users, licenser_users) VALUES (2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)")
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
	ExecCommand("ALTER TABLE licenser ADD " + action + " BOOL")
	ExecCommand("UPDATE licenser SET " + action + " = 1 WHERE id = 0")
	ExecCommand("UPDATE licenser SET " + action + " = 0 WHERE id != 0")
}

// AddDefaultPermissions adds 3 default permissions
func AddDefaultPermissions(name string) {
	AddPermission("viewer_" + name)
	AddPermission("editor_" + name)
	AddPermission("licenser_" + name)
}
