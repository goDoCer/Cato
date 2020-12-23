package parser

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func init() {
	if err := getAuth(); err != nil {
		panic(err)
	}
}

const (
	cateURL      = "https://cate.doc.ic.ac.uk"
	timeTableURL = cateURL + "/timetable.cgi?keyt=%s:%s:%s:%s"
)

var auth string

func login(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func getAuth() error {
	var s map[string]string
	file, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return err
	}
	auth = "Basic " + s["Auth"]
	return nil
}
