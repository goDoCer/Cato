package cate

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/* This file contains functions to convert between "cate time" and real time.
 */

const hoursInADay = 24

var (
	//This variable is used to cache the results of getTermStart() when it is called
	//It will only be set during a cate fetch
	termStart time.Time
)

func getAcademicYear() (currentYear int) {
	//currentYear is the year of last September
	now := time.Now()
	currentYear = now.Year() - 1
	if now.After(time.Date(now.Year(), time.September, 1, 0, 0, 0, 0,
		now.Location())) {
		currentYear++
	}
	return currentYear
}

//Finds the first month of the term, then finds the first day in the table
//The start date is (day - 2) / month / year
func getTermStart(doc *goquery.Document) time.Time {
	var month time.Month
	tableHeaders := doc.Find("th")
	tableHeaders.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		maybeMonth, err := time.Parse("January", strings.TrimSpace(s.Text()))
		if err != nil {
			return true
		}
		month = maybeMonth.Month()
		return false
	})

	var day int
	tableHeaders.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		maybeDay, err := strconv.Atoi(s.Text())
		if err != nil {
			return true
		}
		day = maybeDay
		return false
	})
	t := time.Date(time.Now().Year(), month, day, 0, 0, 0, 0, time.Local)
	return t.AddDate(0, 0, -2)
}

func convertDaysToDate(days int) time.Time {
	return termStart.AddDate(0, 0, days)
}

func convertDateToDays(date time.Time) int {
	return int(date.Sub(termStart).Hours()) / hoursInADay
}

func isWeekend() bool {
	day := time.Now().Weekday()
	return day == time.Saturday || day == time.Sunday
}
