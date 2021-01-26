package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

//getLoginDetails asks the user for their login details
func getLoginDetails() string {
	auth, err := getAuth()
	if err != nil {
		fmt.Println("Enter your shortcode")
		reader := bufio.NewReader(os.Stdin)
		shortcode, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		fmt.Println("Enter your password")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		str := base64.StdEncoding.EncodeToString([]byte(string(shortcode) + ":" + string(bytePassword)))
		return "Basic " + str
	}
	return auth
}

func getAuth() (string, error) {
	var s map[string]string
	file, err := ioutil.ReadFile(path + "/" + "secrets.json")
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return "", err
	}
	return "Basic " + s["Auth"], nil
}
