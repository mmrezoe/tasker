package main

import (
	"fmt"
)

func main() {
	e := fmt.Errorf("test")

	// Compare the string representation of the errors
	if e != nil && e.Error() != fmt.Errorf("test").Error() {
		fmt.Println("Error doesn't match")
		return
	} else {
		fmt.Println("Error matched")
	}
}
