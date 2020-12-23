package cate

import (
	"fmt"
	"log"
	"os"
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

//Fetch clears the stored cache and replaces it with up to date information
func Fetch() {
	fetchInfo()
	fetchModules()
}

//DownloadModule downloads every file in a module
func DownloadModule(module *Module) {
	//Make sure the location exists
	err := checkDir("files/" + module.Name)
	if err != nil {
		log.Fatalln("Failed to create directory", "files/"+module.Name)
	}
	err = downloadModule(module)
	if err != nil {
		fmt.Println(err)
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
	getModules(doc)
	storeModules()
}

func checkDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		return err
	}
	return nil
}
