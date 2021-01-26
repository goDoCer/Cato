package website

import (
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

//This file contains functions related to downloading cate webpages

var extractFilename = regexp.MustCompile("filename=\"(.*)\"")

//GetPage downloads a page from a url
//TODO figure out if the page has given an error
func GetPage(url string) (*goquery.Document, error) {
	resp, err := login(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return goquery.NewDocumentFromResponse(resp)
}

//GetFile downloads a file from cate
func GetFile(url string) (string, []byte, error) {
	resp, err := login(url)
	if err != nil {
		return "", []byte{}, err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", []byte{}, err
	}
	contentHeader := resp.Header.Get("Content-Disposition")
	filename := extractFilename.FindAllStringSubmatch(contentHeader, 1)[0][1]
	if filename == "" {
		return "", []byte{}, errors.New("No filename found")
	}
	return filename, html, err
}
