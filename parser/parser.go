package parser

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func loadFile() (*goquery.Document, error) {
	f, err := os.Open(pageLocation)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return goquery.NewDocumentFromReader(f)
}

func Init() {
	doc, err := loadFile()
	if err != nil {
		log.Fatal(err)
	}
	err = info.getYearAndCourse(doc)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(*info)
}

func getTables() error {
	doc, err := loadFile()
	if err != nil {
		return err
	}
	doc.Find("[src='icons/arrowredright.gif']").Each(func(i int, _ *goquery.Selection) {
		fmt.Println(i)
	})
	return nil
}
