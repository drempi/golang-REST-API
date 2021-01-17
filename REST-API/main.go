package main

import (
	cryptpack "daniel/project/cryptPack"
	databasepack "daniel/project/databasePack"
	httppack "daniel/project/httpPack"
	querypack "daniel/project/queryPack"
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
	app.Get("/", httppack.HelloWorld)
	app.Get("/logout", httppack.LogoutPage)
	app.Get("/test:val", httppack.TestPage)
	app.Post("/", httppack.HelloWorld)
}

func visit(c *fiber.Ctx) error {
	fmt.Println("A visit has been spotted")
	var success bool
	success, querypack.INFO = querypack.StringToQuery(c.Cookies("Visit"))
	// If certain conditions are not met, the user is not logged in.
	if !success || !querypack.INFO.LoggedIn || querypack.INFO.Time.Before(time.Now()) {
		httppack.LoginPage(c)
		querypack.AddCookie(c)
		if querypack.INFO.LoggedIn {
			// Here you continue with your http request stored in cookie. The method is GET.
			c.Method("GET")
			c.Path(querypack.INFO.URL)
			return c.Next()
		}
		// You are not logged in. Appropiate message regarding your request
		// should have shown up somewhere earlier in the program in the LoginPage function
		return nil
	}
	/*
		There should be a list of sites when you DO NOT
		overwrite the URL in the cookie. So far there is
		only one: "/logout". Another solution would be
		to just add this single line whenever necessary
	*/
	if c.OriginalURL() != "/logout" {
		querypack.INFO.URL = c.OriginalURL()
	}

	querypack.AddCookie(c)

	return c.Next()
}
