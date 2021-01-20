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
	loadInfo()
	loadModules()

	err := checkDir(path + "/" + "files")
	if err != nil {
		log.Fatalln("Failed to create files directory")
	}
}

//Login allows a user to store their login details
func Login() error {
	data, _ := json.Marshal(map[string]string{
		"Auth": strings.Replace(auth, "Basic", "", 1),
	})
	err := ioutil.WriteFile(path+"/"+"secrets.json", data, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't save login details\n%s", err.Error())
	}
	return nil
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
func DownloadTask(task *Task, mod *Module) error {
	if task.Downloaded {
		return nil
	}
	//Make sure the location exists
	dir := ModulePath(mod)
	err := checkDir(dir)
	if err != nil {
		return err
	}
	defer storeModules()
	for i, link := range task.Links {
		filename, err := downloadFile(cateURL+"/"+link, dir)
		if err != nil {
			return fmt.Errorf("Couldn't download %s", filename)
		}
		task.FileNames[i] = filename
	}
	task.Downloaded = true
	return nil
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
	fmt.Println(auth)
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
