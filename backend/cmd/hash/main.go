package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "password to hash")
	flag.Parse()

	if *password == "" {
		fmt.Println("Usage: hash -password 'your_password'")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(string(hash))
}
