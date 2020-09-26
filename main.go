package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	s, err := loadSecrets()
	if err != nil {
		panic(err)
	}
	auth := getAuth(s)
	resp, err := login(auth)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(html))
}
