package cate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
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
	modulePath = "modules.json"
)

var (
	coloursToGroups = map[string]int{
		"#f0ccf0": Group,
		"#ccffcc": Individual,
		"white":   Unassessed,
		"#cdcdcd": UnassessedSub,
	}
	Modules = make([]*Module, 0)
)

//Module represents a module for example Reasoning about Programs
type Module struct {
	Name  string
	Tasks []*Task
}

//Task represents a block in the cate timetable
type Task struct {
	Name     string
	Class    int
	Deadline int
	Files    []string //The links to notes for the task
}

//getModules auto loads the modules in a term and all of their current tasks
func getModules(doc *goquery.Document) {
	doc.Find("[style='border: 2px solid blue']").Each(
		func(_ int, sel *goquery.Selection) {
			Modules = append(Modules, parseModule(sel))
		},
	)
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
		Name:  sel.Find("b").Text(),
		Tasks: tasks,
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
			if exists && !strings.Contains(file, "mailto") && !strings.Contains(file, "handins") {
				if strings.Contains(file, "given") {
					files = append(files, getGivenFiles(file)...)
				} else {
					files = append(files, file)
				}
			}
		},
	)
	//Remove redundant whitespace
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(sel.Text(), " ")
	return &Task{
		Name:     strings.TrimSpace(s),
		Class:    coloursToGroups[colour],
		Deadline: day,
		Files:    files,
	}
}

//Some files are "given" and so are on a different page than others
//This function parses that page and returns all the files on it
func getGivenFiles(url string) []string {
	url = cateURL + "/" + url
	data, err := get(url)
	if err != nil {
		log.Println("Error retrieving file from url:", url)
		return []string{}
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		log.Println("Couldn't parse file from url:", url)
		return []string{}
	}
	files := make([]string, 0)
	doc.Find("[href]").Each(func(_ int, sel *goquery.Selection) {
		file, _ := sel.Attr("href")
		if strings.Contains(file, "showfile") {
			files = append(files, file)
		}
	})
	return files
}

//stores the module struct in modules.json
func storeModules() {
	data, err := json.MarshalIndent(Modules, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(modulePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadModules() error {
	data, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Modules)
}
