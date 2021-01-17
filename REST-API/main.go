package main

import (
	cryptpack "github.com/REST-API/cryptPack"
	databasepack "github.com/REST-API/databasePack"
	httppack "github.com/REST-API/httpPack"
	querypack "github.com/REST-API/queryPack"

	"fmt"
	"time"

	"github.com/gofiber/fiber"
	_ "github.com/mattn/go-sqlite3"
)

// DEBUG am i debugging?
var DEBUG bool = false

func debug(s string) {
	fmt.Println("-------------    " + s + "    -------------")
	for i := 0; i < len(s); i++ {
		fmt.Print("=")
	}
	fmt.Println("====================================")
}

func main() {
	databasepack.Init()
	defer databasepack.D.Close()
	cryptpack.InitializeCrypt()

	app := fiber.New()

	app.Use(visit)
	setupRoutes(app)
	app.Listen(":3000")
}

func setupRoutes(app *fiber.App) {
	app.Get("/page", httppack.HelloWorld)
	app.Get("/page/test:val", httppack.TestPage)

	app.Post("/log/login", httppack.LoginPage)
	app.Get("/log/logout", httppack.LogoutPage)
	app.Post("/log/remove", httppack.RemovePage)
	app.Get("/log", httppack.LogPage)
	app.Get("/page", httppack.HelloWorld)
}

func visit(c *fiber.Ctx) error {
	fmt.Println("A visit has been spotted")

	var success bool
	success, querypack.INFO = querypack.StringToQuery(c.Cookies("Visit"))

	// If certain conditions are not met, the user is moved to the logging in site.
	if !success || querypack.INFO.Time.Before(time.Now()) {
		c.Method("GET")
		c.Path("/log")
		querypack.INFO.LoggedIn = false
	}

	if len(c.OriginalURL()) >= 5 && c.OriginalURL()[0:5] == "/page" {
		// If URL has prefix /page then save it in the cookie.
		querypack.INFO.URL = c.OriginalURL()
		if !querypack.INFO.LoggedIn {
			c.Path("/log")
			c.Method("GET")

		}
	} else if !success {
		querypack.INFO.URL = "/page"
	}

	querypack.AddCookie(c)

	return c.Next()
}
