package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//An enum representing the types of tasks
const (
	Group = iota + 1
	Individual
	Unassessed
	UnassessedSub
)

var coloursToGroups = map[string]int{
	"#f0ccf0": Group,
	"#ccffcc": Individual,
	"white":   Unassessed,
	"#cdcdcd": UnassessedSub,
}

//Module represents a module for example Reasoning about Programs
type Module struct {
	name  string
	tasks []*Task
}

//Task represents a block in the cate timetable
type Task struct {
	name     string
	class    int
	deadline int
	files    []string //The links to notes for the task
}

//GetModules auto loads the modules in a term and all of their current tasks
func GetModules(doc *goquery.Document) []*Module {
	modules := make([]*Module, 0)
	doc.Find("[style='border: 2px solid blue']").Each(
		func(_ int, sel *goquery.Selection) {
			modules = append(modules, parseModule(sel))
		},
	)
	return modules
}

func parseModule(sel *goquery.Selection) *Module {
	day := 0
	tasks := make([]*Task, 0)
	sel.Parent().Find("[colspan]").Each(
		func(_ int, sel *goquery.Selection) {
			days, _ := strconv.Atoi(sel.AttrOr("colspan", "0"))
			day += days
			task := parseTask(sel, day)
			if task != nil {
				tasks = append(tasks, task)
			}
		},
	)
	return &Module{
		name:  sel.Find("b").Text(),
		tasks: tasks,
	}
}

func parseTask(sel *goquery.Selection, day int) *Task {
	//Filter out all blank space
	colour, exists := sel.Attr("bgcolor")
	if !exists {
		return nil
	}
	files := make([]string, 0)
	sel.Find("[href]").Each(
		func(_ int, sel *goquery.Selection) {
			file, exists := sel.Attr("href")
			if exists && !strings.Contains(file, "mailto") {
				files = append(files, file)
			}
		},
	)
	//Remove redundant whitespace
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(sel.Text(), " ")
	return &Task{
		name:     strings.TrimSpace(s),
		class:    coloursToGroups[colour],
		deadline: day,
		files:    files,
	}
}

// func L() *goquery.Document {
// 	doc, _ := loadFile("table.html")
// 	return doc
// }
