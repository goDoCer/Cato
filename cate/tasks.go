package cate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	//Modules represents all the courses in the current term
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
	Deadline string   //A string representing time
	Files    []string //The links to notes for the task
}

func findModule(mod *Module) (int, error) {
	for i, m := range Modules {
		if m.Name == mod.Name {
			return i, nil
		}
	}
	return 0, errors.New("Module not found")
}

//parseModules loads the modules in a term and all of their current tasks
//Any preexisting modules are replaced
func parseModules(doc *goquery.Document) {
	doc.Find("[style='border: 2px solid blue']").Each(
		func(_ int, sel *goquery.Selection) {
			mod := parseModule(sel)
			if i, err := findModule(mod); err != nil {
				Modules = append(Modules, mod)
			} else {
				Modules[i] = mod
			}
		},
	)
}

func parseModule(sel *goquery.Selection) *Module {
	today := convertDateToDays(time.Now())
	weekend := isWeekend()
	day := 0
	tasks := make([]*Task, 0)
	sel.Parent().Find("[colspan]").Each(
		func(_ int, sel *goquery.Selection) {
			days, _ := strconv.Atoi(sel.AttrOr("colspan", "0"))
			//Filter out all blank space
			colour, exists := sel.Attr("bgcolor")
			day += days
			if !exists {
				return
			}
			//This accounts for the "blue shift" on weekends
			if weekend && day-days < today {
				day--
			}
			task := parseTask(sel, day, colour)
			tasks = append(tasks, task)
		},
	)
	return &Module{
		Name:  sel.Find("b").Text(),
		Tasks: tasks,
	}
}

func parseTask(sel *goquery.Selection, day int, colour string) *Task {
	files := make([]string, 0)
	//Search all href tags to find links the point to files
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
		Deadline: convertDaysToDate(day).Format(time.ANSIC),
		Files:    files,
	}
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
