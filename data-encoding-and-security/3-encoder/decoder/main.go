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

	fileName := filepath.Join(destDir, "user.json")

	//Open the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// decode the file
	var u user

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&u); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", u)

}
