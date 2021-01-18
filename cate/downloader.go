package cate

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/* This file contains functions related to downloading cate webpages
 */

var extractFilename = regexp.MustCompile("filename=\"(.*)\"")

func get(url string) ([]byte, error) {
	resp, err := login(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return html, nil
}

func downloadFile(url, location string) (string, error) {
	resp, err := login(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	contentHeader := resp.Header.Get("Content-Disposition")
	filename := extractFilename.FindAllStringSubmatch(contentHeader, 1)[0][1]
	if filename == "" {
		fmt.Println(resp.Header)
		return "", errors.New("No filename found")
	}
	if err != nil {
		return "", err
	}
	return filename, ioutil.WriteFile(location+"/"+filename, html, 0644)
}

func downloadHome() (*goquery.Document, error) {
	home, err := get(cateURL)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewBuffer(home))
}

//DownloadTimeTable needs info to be initialised before being called
func downloadTimeTable() (*goquery.Document, error) {
	currentYear := getAcademicYear()

	timetable, err := get(fmt.Sprintf(timeTableURL, currentYear, info.Term,
		info.Code, info.Shortcode))
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewBuffer(timetable))
}

func download(url, location string) error {
	html, err := get(url)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(location, html, 0644)
}

//Some files are "given" and so are on a different page than others
//This function parses that page and returns all the files on it
func getGivenFiles(url string) []string {
	url = cateURL + "/" + url
	data, err := get(url)
	if err != nil {
		log.Println("Error retrieving file from url:", url)
		return []string{}
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		log.Println("Couldn't parse file from url:", url)
		return []string{}
	}
	files := make([]string, 0)
	doc.Find("[href]").Each(func(_ int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		if strings.Contains(link, "showfile") {
			files = append(files, link)
		}
	})
	return files
}
