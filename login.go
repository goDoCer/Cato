package main

import (
	"net/http"
)

func getAuth(s secrets) string {
	return "Basic " + s["Auth"]
}

func login(auth string) (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://cate.doc.ic.ac.uk", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}
