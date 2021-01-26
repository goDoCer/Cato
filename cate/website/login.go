package website

import (
	"errors"
	"net/http"
)

//URLs related to specific cate webpages
const (
	CateURL      = "https://cate.doc.ic.ac.uk"
	TimeTableURL = CateURL + "/timetable.cgi?keyt=%d:%d:%s:%s"
)

var auth string

//SetAuth should be used to give this package the auth token
func SetAuth(token string) {
	auth = token
}

func login(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Access forbidden")
	}
	return resp, err
}
