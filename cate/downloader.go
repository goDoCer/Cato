package cate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

func download(url, location string) error {
	html, err := get(url)
	if err != nil {
		return err
	}
	fmt.Println(location)
	return ioutil.WriteFile(location, html, 0644)
}

//TODO only save the details needed rather than the entire file
func downloadHome() (*goquery.Document, error) {
	home, err := get(cateURL)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewBuffer(home))
}

//DownloadTimeTable needs info to be initialised before being called
func downloadTimeTable() (*goquery.Document, error) {
	//Current currentYear is the year of last September
	var currentYear int
	now := time.Now()
	currentYear = now.Year() + 1
	if now.After(time.Date(now.Year(), time.September, 1, 0, 0, 0, 0,
		now.Location())) {
		currentYear--
	}

	timetable, err := get(fmt.Sprintf(timeTableURL, currentYear, info.Term,
		info.Code, info.Shortcode))
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewBuffer(timetable))
}

//DownloadModule tries to download all tasks in a module in the appropriate folder
//It stops downloading as soon as it fails once
func downloadModule(module *Module) error {
	var err error
	location := "files/" + formatName(module.Name) + "/"
	for _, task := range module.Tasks {
		for _, file := range task.Files {
			err = download(cateURL+"/"+file, location+formatName(task.Name)+".pdf")
			if err != nil {
				log.Println("Error downloading module: " + module.Name)
				return err
			}
		}
	}
	return nil
}

func formatName(name string) string {
	return strings.ReplaceAll(name, ":", "")
}
