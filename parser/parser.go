package parser

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func loadFile(path string) (*goquery.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return goquery.NewDocumentFromReader(f)
}

//Init initialises all singletons when a new cate.html file is loaded in
func Init() {
	doc, err := loadFile(pageLocation)
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
	fmt.Println(*info)
}
