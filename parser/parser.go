package parser

import (
	"log"
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
}

//Fetch clears the stored cache and replaces it with up to date information
func Fetch() {
	fetchInfo()
	fetchModules()
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
