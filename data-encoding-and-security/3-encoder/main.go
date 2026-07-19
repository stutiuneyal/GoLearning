package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	destDir = "data-encoding-and-security/3-encoder/json"
)

type user struct {
	Name     string `json:"username"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Phone    string `json:"mobile"`
	IsActive bool   `json:"is_active"`
}

func main() {

	u1 := user{
		Name:     "Jane",
		Age:      32,
		Password: "haha",
		Phone:    "1234567890",
		IsActive: true,
	}

	// u2 := user{
	// 	Name:     "John",
	// 	Age:      28,
	// 	Password: "lala",
	// 	Phone:    "1234567880",
	// 	IsActive: false,
	// }

	//encoder := json.NewEncoder(os.Stdout)

	if err := os.MkdirAll(destDir, 0o700); err != nil {
		fmt.Println(err)
	}

	fileName := filepath.Join(destDir, "user.json")

	// file, err := os.Create(fileName)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer file.Close()

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o660)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(u1); err != nil {
		fmt.Println(err)
	}

}
