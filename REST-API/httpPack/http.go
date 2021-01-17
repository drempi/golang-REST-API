package httppack

import (
	querypack "daniel/project/queryPack"
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

// LoginPage its a login site
func LoginPage(c *fiber.Ctx) {
	if querypack.INFO.LoggedIn {
		querypack.INFO.URL = c.OriginalURL()
	}
	querypack.INFO.LoggedIn = false
	if c.Method() == "POST" {
		var acc Account
		err := c.BodyParser(&acc)
		if err != nil {
			fmt.Println(err.Error())
			c.Send([]byte("Wrong input."))
			return
		}
		role := FindUser(acc)
		if acc.Isnew {
			fmt.Println("NEW")
			if role == -2 {
				querypack.INFO.LoggedIn = true
				c.Send([]byte("Created account!"))
				CreateUser(acc)
			} else {
				c.Send([]byte("Account with such login already exists."))
			}
		} else {
			fmt.Println("NOT NEW")
			if role == -1 {
				c.Send([]byte("Wrong password!"))
			} else if role == -2 {
				c.Send([]byte("No such username!"))
			} else {
				querypack.INFO.LoggedIn = true
				c.Send([]byte("Successfully logged in!"))
			}
		}
	} else {
		c.Send([]byte("You are not logged in. Type your username and password or create new account."))
	}
}

// LogoutPage Its when you log out of your account
func LogoutPage(c *fiber.Ctx) error {
	querypack.INFO.LoggedIn = false
	querypack.AddCookie(c)
	return c.Send([]byte("logging out!"))
}

// TestPage its a test page
func TestPage(c *fiber.Ctx) error {
	return c.Send([]byte("test nr " + c.Params("val")))
}
