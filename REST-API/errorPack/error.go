package errorpack

import (
	"fmt"
	"os"
)

// It's a small pack to handle some errors

// OK checks if error is null. If not, terminates the server
func OK(err *error) {
	if (*err) != nil {
		fmt.Println("FATAL ERROR!")
		fmt.Println((*err).Error())
		os.Exit(1)
	}
}

// CHECK checks if error is null. If not, prints it out
func CHECK(err *error) bool {
	if (*err) != nil {
		fmt.Println("ERROR SPOTTED!")
		fmt.Println((*err).Error())
		return true
	}
	return false
}
