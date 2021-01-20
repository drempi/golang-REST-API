package httppack

import (
	databasepack "github.com/drempi/golang-REST-API/REST-API/databasePack"
	querypack "github.com/drempi/golang-REST-API/REST-API/queryPack"

	"fmt"
	"strconv"

	"github.com/gofiber/fiber"
)

// MainPage its a main page
func MainPage(c *fiber.Ctx) error {
	return HelloWorld(c)
}

// HelloWorld it just prints "Hello world"
func HelloWorld(c *fiber.Ctx) error {
	return c.Send([]byte("Hello World! Default Page"))
}

// LoadPage loads the last page after successfully logging in
func LoadPage(c *fiber.Ctx, login string, roles []int) error {
	querypack.INFO.Roles = roles
	querypack.INFO.Login = login
	querypack.INFO.LoggedIn = true
	c.Method("GET")
	c.Path(querypack.INFO.URL)
	querypack.AddCookie(c)
	return c.Next()
}

// LoginPage its a login site
func LoginPage(c *fiber.Ctx) error {
	var acc databasepack.Account
	err := c.BodyParser(&acc)
	if err != nil {
		fmt.Println(err.Error())
		return c.Send([]byte("Wrong input."))
	}
	roles := databasepack.FindUser(acc)
	if acc.Isnew {
		//fmt.Println("NEW")
		if roles[0] == -2 {
			databasepack.CreateUser(acc)
			return LoadPage(c, acc.Login, []int{0})
		}
		return c.Send([]byte("Account with such login already exists."))
	}
	//fmt.Println("NOT NEW")
	if roles[0] == -1 {
		return c.Send([]byte("Wrong password!"))
	} else if roles[0] == -2 {
		return c.Send([]byte("No such username!"))
	}
	return LoadPage(c, acc.Login, roles)
}

// LogoutPage Its when you log out of your account
func LogoutPage(c *fiber.Ctx) error {
	querypack.INFO.LoggedIn = false
	querypack.AddCookie(c)
	return c.Send([]byte("logging out!"))
}

// RemovePage removes user
func RemovePage(c *fiber.Ctx) error {
	var acc databasepack.Account
	err := c.BodyParser(&acc)
	if err != nil {
		fmt.Println(err.Error())
		return c.Send([]byte("Wrong input."))
	}
	roles := databasepack.FindUser(acc)
	if roles[0] < 0 {
		return c.Send([]byte("Bad login or password."))
	}
	databasepack.RemoveUser(acc)
	return LogPage(c)
}

// LogPage Default Login page
func LogPage(c *fiber.Ctx) error {
	return c.Send([]byte("You are on a default Login page."))
}

// TestPage its a test page
func TestPage(c *fiber.Ctx) error {
	return c.Send([]byte("test nr " + c.Params("val")))
}

// AddTablePage adding tables to the database
func AddTablePage(c *fiber.Ctx) error {
	// Firstly checks if the user can actually perform this action
	if !databasepack.Allowed(querypack.INFO.Roles, "creator_") {
		return c.Send([]byte("Permission declined!"))
	}
	// Firstly check if all the informations are correct.
	var tab databasepack.TableType
	err := c.BodyParser(&tab)
	if err != nil || tab.Type < 0 || tab.Type > 1 {
		fmt.Println(err.Error())
		return c.Send([]byte("Wrong input."))
	}
	exists := databasepack.FindTable(tab.Name)
	if exists != -1 {
		return c.Send([]byte("Table with such name already exists!"))
	}
	// Next check if the name is proper
	if !databasepack.CheckName([]byte(tab.Name)) {
		return c.Send([]byte("Bad table name!"))
	}
	databasepack.AddTable(tab)
	return c.Send([]byte("Added table: " + tab.Name))
}

// EditPage this is where you add content to the database
func EditPage(c *fiber.Ctx) error {
	// Firstly checks if given table exists
	name := c.Params("name")
	T := databasepack.FindTable(name)
	if T < 0 {
		return c.Send([]byte("No table with such name!"))
	}
	// Next checks if the user can actually perform this action
	if !databasepack.Allowed(querypack.INFO.Roles, "editor_"+name) {
		return c.Send([]byte("Permission declined!"))
	}
	// Lastly posts given post
	databasepack.CreatePosts(c.Body(), name, T)
	return c.Send([]byte("Successfully edited!"))
}

// GetPage gets elemenents
func GetPage(c *fiber.Ctx) error {
	// First check if offset and amt are integers
	offset, err := strconv.Atoi(c.Params("offset"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	amt, err := strconv.Atoi(c.Params("amt"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	lang, err := strconv.Atoi(c.Params("lang"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	// Then checks if given table exists
	name := c.Params("name")
	T := databasepack.FindTable(name)
	if T < 0 {
		return c.Send([]byte("No table with such name!"))
	}
	// Next checks if the user can view these elements
	if !databasepack.Allowed(querypack.INFO.Roles, "viewer_"+name) {
		return c.Send([]byte("Permission declined!"))
	}

	// Gives the ordered elements
	return c.Send(databasepack.GetPosts(name, T, offset, amt, lang))
}

// LicensePage changes the licenser in given field
func LicensePage(c *fiber.Ctx) error {
	// First check if group and val are integers
	group, err := strconv.Atoi(c.Params("group"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	val, err := strconv.ParseBool(c.Params("val"))
	if err != nil {
		return c.Send([]byte("Value not a bool!"))
	}
	action := c.Params("action")
	// Next checks if the user can view these elements
	if !databasepack.Allowed(querypack.INFO.Roles, "licenser_") {
		return c.Send([]byte("Permission declined!"))
	}
	// changes the value in licenser
	databasepack.EditLicenser(group, action, val)
	return c.Send([]byte("Successfully changed licenser table"))
}

// AddGroupUserPage adds group to some user
func AddGroupUserPage(c *fiber.Ctx) error {
	// First check if group is an integers
	group, err := strconv.Atoi(c.Params("group"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	login := c.Params("login")
	// Next checks if the user can edit table users
	if !databasepack.Allowed(querypack.INFO.Roles, "editor_users") {
		return c.Send([]byte("Permission declined!"))
	}
	// changes groups
	return c.Send([]byte(databasepack.AddGroupUser(login, group)))
}

// RemoveGroupUserPage removes group from some user
func RemoveGroupUserPage(c *fiber.Ctx) error {
	// First check if group is an integers
	group, err := strconv.Atoi(c.Params("group"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	login := c.Params("login")
	// Next checks if the user can edit table users
	if !databasepack.Allowed(querypack.INFO.Roles, "editor_users") {
		return c.Send([]byte("Permission declined!"))
	}
	// changes groups
	return c.Send([]byte(databasepack.RemoveGroupUser(login, group)))
}

// AddGroupPage adds an empty group to the database
func AddGroupPage(c *fiber.Ctx) error {
	// First check if id is an integer
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Send([]byte("Value not an integer!"))
	}
	// Next checks if the user can edit table licenser
	if !databasepack.Allowed(querypack.INFO.Roles, "licenser_") {
		return c.Send([]byte("Permission declined!"))
	}
	return c.Send([]byte(databasepack.AddGroup(id)))
}

// ResetPasswordPage sends the email with link to change the password
func ResetPasswordPage(c *fiber.Ctx) error {
	var acc databasepack.Account
	err := c.BodyParser(&acc)
	if err != nil {
		fmt.Println(err.Error())
		return c.Send([]byte("Wrong input."))
	}
	// The rest is here
	return c.Send([]byte(databasepack.ResetPassword(acc)))
}

// RandomPasswordEmailPage sends the email with link to change the password
func RandomPasswordEmailPage(c *fiber.Ctx) error {
	login := c.Params("login")
	// Try to send the email
	return c.Send([]byte(databasepack.SendEmail(login)))
}

// RandomPasswordPage sets the password of user to something random and outputs it
func RandomPasswordPage(c *fiber.Ctx) error {
	login := c.Params("login")
	hash := c.Params("hash")
	// Just change it, if everything goes okay that is
	return c.Send([]byte(databasepack.RandomPasswordUser(login, hash)))
}
