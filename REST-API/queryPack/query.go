package querypack

import (
	cryptpack "github.com/drempi/golang-REST-API/REST-API/cryptPack"

	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber"
)

// Query contains all the information about a certain http request
type Query struct {
	URL      string    `json:"url"`
	LoggedIn bool      `json:"logged_in"`
	Login    string    `json:"login"`
	Time     time.Time `json:"time"`
	Role     int       `json:"role"`
}

// INFO its all there is!
var INFO Query

// QueryToString converts query to encrypted string
func QueryToString(Q Query) string {
	message, err := json.Marshal(Q)
	if err != nil {
		log.Fatalln("json conversion failed")
	}

	var temp []byte
	success, temp := cryptpack.Encrypt([]byte(message))
	if !success {
		log.Fatalln("could not encrypt!")
	}
	return string(cryptpack.SmallBase(temp))
}

// StringToQuery converts encrypted string to query object
func StringToQuery(s string) (bool, Query) {
	var Q Query
	success, message := cryptpack.Decrypt(cryptpack.BigBase([]byte(s)))
	if !success {
		Q.LoggedIn = true
		return false, Q
	}
	err := json.Unmarshal([]byte(message), &Q)
	if err != nil {
		Q.LoggedIn = true
		return false, Q
	}
	return true, Q
}

// AddCookie adds cookie based on current INFO
func AddCookie(c *fiber.Ctx) {
	INFO.Time = time.Now().Add(time.Hour)
	encrypted := QueryToString(INFO)
	cookie := new(fiber.Cookie)
	//cookie.HTTPOnly = true
	cookie.Name = "Visit"
	cookie.Value = encrypted
	c.Cookie(cookie)
}
