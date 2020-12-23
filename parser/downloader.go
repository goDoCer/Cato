package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	pageLocation      = "cate.html"
	timetableLocation = "table.html"
)

func download(url, location string) error {
	// resp, err := login(url)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()
	// html, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	html := []byte("hello world")
	fmt.Println(location)
	return ioutil.WriteFile(location, html, 0644)
}

//TODO only save the details needed rather than the entire file
func DownloadHome() error {
	return download(cateURL, pageLocation)
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

	return download(fmt.Sprintf(timeTableURL, currentYear, info.term, info.code,
		info.shortcode), timetableLocation)
}

//DownloadModule tries to download all tasks in a module in the appropriate folder
func DownloadModule(module *Module) error {
	var err error
	location := "files/" + formatName(module.name) + "/"
	if _, err := os.Stat(location); os.IsNotExist(err) {
		os.Mkdir(location, os.ModePerm)
	}
	for _, task := range module.tasks {

		for _, file := range task.files {
			err = download(cateURL+"/"+file, location+formatName(task.name)+".pdf")
			if err != nil {
				log.Println("Error downloading module: " + module.name)
				return err
			}
		}
	}
	return nil
}

func formatName(name string) string {
	return strings.ReplaceAll(name, ":", "")
}
