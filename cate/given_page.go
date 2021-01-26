package cate

import (
	"log"
	"strings"

	"github.com/Akshat-Tripathi/cateCli/cate/website"
	"github.com/PuerkitoBio/goquery"
)

// Some files are "given" and so are on a different page than others
// This function parses that page and returns all the files on it
func getGivenFiles(url string) []string {
	doc, err := website.GetPage(website.CateURL + "/" + url)
	if err != nil {
		log.Println("Couldn't parse file from url:", url)
		return []string{}
	}
	files := make([]string, 0)
	doc.Find("[href]").Each(func(_ int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		if strings.Contains(link, "showfile") {
			files = append(files, link)
		}
	})
	return files
}
