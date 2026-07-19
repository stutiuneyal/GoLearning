package main

import (
	"encoding/json"
	"fmt"
)

type user struct {
	Name     string `json:"username"`
	Age      int
	Password string `json:"-"`
	Phone    string `json:"mobile"`
	IsActive bool   `json:"is_active"`
	Role     string
	Profile  Profile `json:"profile"`
}

type Profile struct {
	URL string `json:"url"`
}

var payload = `
{
	"username":"Jane",
	"age": 32,
	"password":"1234",
	"mobile":"1234-5678-90",
	"is_active": true,
	"profile": {
		"url": "https://google.com"
	}
}
`

func main() {

	var jane user

	if err := json.Unmarshal([]byte(payload), &jane); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("User: %#v", jane)

}
