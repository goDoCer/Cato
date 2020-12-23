package parser

import (
	"io/ioutil"
)

const pageLocation = "cate.html"

//TODO only save the details needed rather than the entire file
func Download() error {
	resp, err := login(getAuth(s), cateURL)
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

func downloadTimeTable() error {

}
