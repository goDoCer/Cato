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

	err = checkDir("files")
	if err != nil {
		log.Fatalln("Failed to create files directory")
	}
}

//Login allows a user to store their login details
func Login() {
	data, _ := json.Marshal(map[string]string{
		"Auth": strings.Replace(auth, "Basic", "", 1),
	})
	err := ioutil.WriteFile("secrets.json", data, 0644)
	if err != nil {
		fmt.Println("Could not save login details because", err)
	}
}

//Fetch clears the stored cache and replaces it with up to date information
func Fetch() {
	fetchInfo()
	fetchModules()
}

//DownloadModule downloads every file in a module
//It stops downloading as soon as it fails once
func DownloadModule(module *Module) {
	//Make sure the location exists
	err := checkDir("files/" + module.Name)
	if err != nil {
		log.Fatalln("Failed to create directory", "files/"+module.Name)
	}
	defer storeModules()
	location := "files/" + strings.ReplaceAll(module.Name, ":", "") + "/"
	for _, task := range module.Tasks {
		if task.Downloaded {
			continue
		}
		for _, link := range task.Files {
			err = downloadFile(cateURL+"/"+link, location)
			if err != nil {
				fmt.Println("Error downloading module: "+module.Name, err)
				return
			}
		}
		task.Downloaded = true
	}

}

func fetchInfo() {
	doc, err := downloadHome()
	if err != nil {
		log.Fatal(err)
	}
	err = getYearAndCourse(doc)
	if err != nil {
		log.Println(err)
	}
	getName(doc)
	getTerm(doc)
	getShortcode(doc)
	storeInfo()
}

func fetchModules() {
	doc, err := downloadTimeTable()
	if err != nil {
		log.Fatal(err)
	}
	//TODO figure out if the page has given an error
	termStart = getTermStart(doc)
	parseModules(doc)
	storeModules()
}

func checkDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		return err
	}
	return nil
}
