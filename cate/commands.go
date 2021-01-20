package cate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Init initialises all singletons when a new cate.html file is loaded in
func Init() {
	err := loadInfo()
	if err != nil {
		fetchInfo()
	}

	err = loadModules()
	if err != nil {
		fetchModules()
	}

	err = checkDir(path + "/" + "files")
	if err != nil {
		log.Fatalln("Failed to create files directory")
	}
}

//Login allows a user to store their login details
func Login() {
	data, _ := json.Marshal(map[string]string{
		"Auth": strings.Replace(auth, "Basic", "", 1),
	})
	err := ioutil.WriteFile(path+"/"+"secrets.json", data, 0644)
	if err != nil {
		fmt.Println("Could not save login details because", err)
	}
}

//Fetch clears the stored cache and replaces it with up to date information
func Fetch() error {
	err := fetchInfo()
	if err != nil {
		return err
	}
	return fetchModules()
}

//DownloadTask downloads a task from a module
func DownloadTask(task *Task, mod *Module) {
	if task.Downloaded {
		return
	}
	//Make sure the location exists
	dir := ModulePath(mod)
	err := checkDir(dir)
	if err != nil {
		log.Fatalln("Failed to create directory", dir, err)
	}
	defer storeModules()
	for i, link := range task.Links {
		filename, err := downloadFile(cateURL+"/"+link, dir)
		if err != nil {
			fmt.Println("Error downloading module: "+mod.Name, err)
			return
		}
		task.FileNames[i] = filename
	}
	task.Downloaded = true
}

//ModulePath returns the path to a module
func ModulePath(mod *Module) string {
	return path + "/files/" + strings.ReplaceAll(mod.Name, ":", "")
}

func fetchInfo() error {
	doc, err := downloadHome()
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
	storeInfo()
	return nil
}

func fetchModules() error {
	doc, err := downloadTimeTable()
	if err != nil {
		return err
	}
	//TODO figure out if the page has given an error
	termStart = getTermStart(doc)
	parseModules(doc)
	return storeModules()
}

func checkDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		return err
	}
	return nil
}
