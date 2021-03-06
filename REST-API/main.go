package main

import (
	cryptpack "github.com/drempi/golang-REST-API/REST-API/cryptPack"
	databasepack "github.com/drempi/golang-REST-API/REST-API/databasePack"
	httppack "github.com/drempi/golang-REST-API/REST-API/httpPack"
	querypack "github.com/drempi/golang-REST-API/REST-API/queryPack"

	"fmt"
	"time"

	"github.com/gofiber/fiber"
	_ "github.com/mattn/go-sqlite3"
)

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
	app.Get("/page/get/:name/:lang/:offset/:amt", httppack.GetPage)
	app.Post("/page/addtable", httppack.AddTablePage)
	app.Post("/page/edit/:name", httppack.EditPage)
	app.Get("/page/licenser/:group/:action/:val", httppack.LicensePage)
	app.Get("/page/addgroupuser/:login/:group", httppack.AddGroupUserPage)
	app.Get("/page/removegroupuser/:login/:group", httppack.RemoveGroupUserPage)
	app.Get("/page/addgroup/:id", httppack.AddGroupPage)

	app.Get("/log/changepassword/:login", httppack.ResetPasswordPage)
	app.Get("/log/resetpasswordemail/:login", httppack.RandomPasswordEmailPage)
	app.Get("/log/resetpassword/:login/:hash", httppack.RandomPasswordPage)
	app.Post("/log/login", httppack.LoginPage)
	app.Get("/log/logout", httppack.LogoutPage)
	app.Post("/log/remove", httppack.RemovePage)
	app.Get("/log", httppack.LogPage)
	app.Get("/page", httppack.HelloWorld)
}

func visit(c *fiber.Ctx) error {
	fmt.Println("A visit has been spotted " + c.Method() + "   " + c.OriginalURL())

	// If user wants to change the password, don't check other things
	if len(c.OriginalURL()) >= 20 && c.OriginalURL()[0:21] == "/log/changepassword/" {
		return c.Next()
	}

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
