package cate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/* This file stores the details of a student i.e. relevant information found on
   the cate homepage.
   See the *details* struct for the information stored.
*/

//info is a singleton containing the details of the user
var info = new(details)

type details struct {
	Name      string
	Shortcode string
	Code      string
	Undergrad bool
	Year      int
	Course    int
	Term      int
}

//enum for the different course types
const (
	comp = iota + 1
	jmc
	ise
	occ
	ai
	compSpec
	hipeds
	research
	industrial
	bio
)

func loadInfo() error {
	data, err := ioutil.ReadFile("info.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &info)
}

func storeInfo() {
	data, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("info.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getName(doc *goquery.Document) {
	info.Name = doc.Find("[style='padding-left: 5px; text-align: left;'][colspan='3']").
		First().
		Children().
		First().
		Text()
}

func getShortcode(doc *goquery.Document) {
	info.Shortcode = strings.Replace(doc.Find("title").Text(), "CATe - ", "", 1)
}

func getTerm(doc *goquery.Document) {
	doc.Find("[name='period']").EachWithBreak(func(_ int, sel *goquery.Selection) bool {
		if _, ok := sel.Attr("checked"); ok {
			term, ok := sel.Attr("value")
			if ok {
				termInt, err := strconv.Atoi(term)
				if err != nil {
					log.Println("Term is not an int")
				}
				info.Term = termInt
			}
			return false
		}
		return true
	})
}

func getYearAndCourse(doc *goquery.Document) error {
	err := errors.New("No class or course information found\nTry cate fetch")
	doc.Find("[name='class']").
		EachWithBreak(func(_ int, sel *goquery.Selection) bool {
			if _, ok := sel.Attr("style"); ok {
				class, ok := sel.Attr("value")
				if ok {
					info.Code = class
					info.Undergrad, info.Year, info.Course, err = parseClassCode(class)
				}
				return false
			}
			return true
		})
	return err
}

//PRE: code is a 2 character string
func parseClassCode(code string) (undergrad bool, year int, course int, err error) {
	year, err = strconv.Atoi(code[1:])
	if err != nil {
		return true, 1, 1, err
	}
	//This should really be refactored into a map loaded from a json file once a settings format is decided
	course, undergrad = func() (int, bool) {
		switch code[0] {
		case 'c':
			return comp, true
		case 'j':
			return jmc, true
		case 'i':
			return ise, true
		case 'o':
			return occ, true
		case 'v':
			return comp, false
		case 't':
			return ai, false
		case 's':
			return compSpec, false
		case 'h':
			return hipeds, false
		case 'r':
			return research, false
		case 'y':
			return industrial, false
		case 'b':
			return bio, false
		default:
			return 0, false
		}
	}()
	if course == 0 {
		return true, 1, 1, errors.New("Invalid course code")
	}
	return undergrad, year, course, nil
}
