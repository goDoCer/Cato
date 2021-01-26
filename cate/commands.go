package cate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Akshat-Tripathi/cateCli/cate/website"
)

// Init initialises all singletons when a new cate.html file is loaded in
func Init(path string) {
	loadInfo(path)
	loadModules(path)

	err := checkDir(path + "/" + "files")
	if err != nil {
		log.Fatalln("Failed to create files directory")
	}
}

//SetAuth allows this module to make requests to cate
func SetAuth(auth string) {
	website.SetAuth(auth)
}

//Fetch clears the stored cache and replaces it with up to date information
func Fetch(path string) error {
	err := fetchInfo(path)
	if err != nil {
		return err
	}
	return fetchModules(path)
}

//DownloadTask downloads a task from a module
func DownloadTask(task *Task, mod *Module, path string) error {
	if task.Downloaded {
		return nil
	}
	//Make sure the location exists
	dir := path + ModulePath(mod)
	err := checkDir(dir)
	if err != nil {
		return err
	}
	defer storeModules(path)
	task.FileNames = make([]string, len(task.Links))
	for i, link := range task.Links {
		filename, data, err := website.GetFile(website.CateURL + "/" + link)
		if err != nil {
			return fmt.Errorf("Couldn't download %s", filename)
		}
		err = ioutil.WriteFile(dir+"/"+filename, data, 0644)
		if err != nil {
			return err
		}
		task.FileNames[i] = filename
	}
	task.Downloaded = true
	return nil
}

//ModulePath returns the relative path to a module
func ModulePath(mod *Module) string {
	return "/files/" + strings.ReplaceAll(mod.Name, ":", "")
}

func fetchInfo(path string) error {
	doc, err := website.GetPage(website.CateURL)
	if err != nil {
		return err
	}
	err = getYearAndCourse(doc)
	if err != nil {
		return err
	}
	getName(doc)
	getTerm(doc)
	getShortcode(doc)
	storeInfo(path)
	return nil
}

func fetchModules(path string) error {
	doc, err := website.GetPage(fmt.Sprintf(website.TimeTableURL,
		getAcademicYear(), info.Term, info.Code, info.Shortcode))
	if err != nil {
		return err
	}
	//TODO figure out if the page has given an error
	termStart = getTermStart(doc)
	parseModules(doc)
	return storeModules(path)
}

func checkDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		return err
	}
	return nil
}
