package parser

import (
	"fmt"
	"io/ioutil"
	"time"
)

const (
	pageLocation      = "cate.html"
	timetableLocation = "table.html"
)

//TODO only save the details needed rather than the entire file
func Download() error {
	resp, err := login(cateURL)
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

//DownloadTimeTable needs info to be initialised before being called
func DownloadTimeTable() error {
	//Current currentYear is the currentYear of last September
	var currentYear int
	now := time.Now()
	currentYear = now.Year() + 1
	if now.After(time.Date(now.Year(), time.September, 1, 0, 0, 0, 0,
		now.Location())) {
		currentYear--
	}

	resp, err := login(fmt.Sprintf(timeTableURL, currentYear, info.term,
		info.code, info.shortcode))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(timetableLocation, html, 0644)
}
