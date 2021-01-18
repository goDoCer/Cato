package cate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/* This file contains functions related to downloading and parsing cate webpages
 */

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

func formatName(name string) string {
	return strings.ReplaceAll(name, ":", "")
}

//Some files are "given" and so are on a different page than others
//This function parses that page and returns all the files on it
func getGivenFiles(url string) map[string]string {
	url = cateURL + "/" + url
	data, err := get(url)
	files := make(map[string]string)
	if err != nil {
		log.Println("Error retrieving file from url:", url)
		return files
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		log.Println("Couldn't parse file from url:", url)
		return files
	}
	doc.Find("[href]").Each(func(_ int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		if strings.Contains(link, "showfile") {
			filename := sel.Text()
			files[filename] = link
		}
	})
	return files
}
