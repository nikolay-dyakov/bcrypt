package main

import (
	"fmt"
	"flag"
	"golang.org/x/crypto/bcrypt"
	"github.com/atotto/clipboard"
)

func main() {
	password := flag.String("p", "bcrypt", "Password to hash")
	cost := flag.Int("c", 10, "Cost")
	copy := flag.Bool("cc", false, "Copy bcrypt hash to clipboard")
	flag.Parse()
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), *cost)
	if err != nil {
		fmt.Println("Sorry cant generate hash for that password")
	} else if *copy {
		clipboard.WriteAll(string(hash))
	} else {
		fmt.Println(string(hash))
	}
}
