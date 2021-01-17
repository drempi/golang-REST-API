package httppack

import (
	querypack "github.com/REST-API/queryPack"

	"fmt"

	"github.com/gofiber/fiber"
)

// MainPage its a main page
func MainPage(c *fiber.Ctx) error {
	return HelloWorld(c)
}

// HelloWorld it just prints "Hello world"
func HelloWorld(c *fiber.Ctx) error {
	return c.Send([]byte("Hello World"))
}

// LoadPage loads the last page after successfully logging in
func LoadPage(c *fiber.Ctx) error {
	querypack.INFO.LoggedIn = true
	c.Method("GET")
	c.Path(querypack.INFO.URL)
	querypack.AddCookie(c)
	return c.Next()
}

// LoginPage its a login site
func LoginPage(c *fiber.Ctx) error {
	var acc Account
	err := c.BodyParser(&acc)
	if err != nil {
		fmt.Println(err.Error())
		return c.Send([]byte("Wrong input."))
	}
	role := FindUser(acc)
	if acc.Isnew {
		//fmt.Println("NEW")
		if role == -2 {
			CreateUser(acc)
			return LoadPage(c)
		}
		return c.Send([]byte("Account with such login already exists."))
	}
	//fmt.Println("NOT NEW")
	if role == -1 {
		return c.Send([]byte("Wrong password!"))
	} else if role == -2 {
		return c.Send([]byte("No such username!"))
	}
	return LoadPage(c)
}

// LogoutPage Its when you log out of your account
func LogoutPage(c *fiber.Ctx) error {
	querypack.INFO.LoggedIn = false
	querypack.AddCookie(c)
	return c.Send([]byte("logging out!"))
}

// RemovePage removes user
func RemovePage(c *fiber.Ctx) error {
	var acc Account
	err := c.BodyParser(&acc)
	if err != nil {
		fmt.Println(err.Error())
		return c.Send([]byte("Wrong input."))
	}
	role := FindUser(acc)
	if role < 0 {
		return c.Send([]byte("Bad login or password."))
	}
	RemoveUser(acc)
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
