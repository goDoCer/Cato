package parser

import (
	"io/ioutil"
)

const pageLocation = "cate.html"

//TODO only save the details needed rather than the entire file
func Download() error {
	s, err := loadSecrets()
	if err != nil {
		return err
	}
	auth := getAuth(s)
	resp, err := login(auth)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pageLocation, html, 0644)
}
