package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	byte, _ := bcrypt.GenerateFromPassword([]byte("passw0rd"), bcrypt.DefaultCost)

	fmt.Println(string(byte))
}
