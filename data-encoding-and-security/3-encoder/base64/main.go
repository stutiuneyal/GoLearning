package main

import (
	"encoding/base64"
	"fmt"
)

func main() {

	data := "Welcome to the world of Go!"

	encodedString := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(encodedString)

	byteData, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byteData))

}
