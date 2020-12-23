package parser

import (
	"log"
)

// Init initialises all singletons when a new cate.html file is loaded in
func Init() {
	err := loadInfo()
	if err == nil {
		return
	}
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
