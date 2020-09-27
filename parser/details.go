package parser

import (
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

//info is a singleton containing the details of the user
var info = new(details)

type details struct {
	name      string
	undergrad bool
	year      int
	course    int
}

//enum for the different courseTypes
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

func (d *details) getName(doc *goquery.Document) {
	d.name = doc.Find("[style='padding-left: 5px; text-align: left;'][colspan='3']").
		First().
		Children().
		First().
		Text()
}

func (d *details) getYearAndCourse(doc *goquery.Document) error {
	err := errors.New("No class or course information found\nTry cate fetch")
	doc.Find("[name='class']").
		EachWithBreak(func(_ int, sel *goquery.Selection) bool {
			if _, ok := sel.Attr("checked"); ok {
				class, ok := sel.Attr("value")
				if ok {
					d.undergrad, d.year, d.course, err = parseClassCode(class)
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
