package cate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
	"sort"
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
	Name       string
	Class      int
	Downloaded bool
	Deadline   time.Time
	Links      []string //Contains links to all associated files
	FileNames  []string
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

//Given a selection containing a module name (in the blue boxes), this finds
//all tasks within the module
func parseModule(sel *goquery.Selection) *Module {
	today := convertDateToDays(time.Now())
	tasks := make([]*Task, 0)
	var day int
	//Get number of rows in module
	nRows, _ := strconv.Atoi(sel.AttrOr("rowspan", "1"))
	row := sel.Parent()
	for i := 0; i < nRows; i++ {
		day = 0
		row.Find("[colspan]").Each(
			func(_ int, sel *goquery.Selection) {
				days, _ := strconv.Atoi(sel.AttrOr("colspan", "0"))
				//Filter out all blank space
				colour, exists := sel.Attr("bgcolor")
				day += days
				if !exists {
					return
				}
				//This accounts for the "blue shift" for tasks set before the current day
				if day-days < today {
					day--
				}
				task := parseTask(sel, day, colour)
				tasks = append(tasks, task)
			},
		)
		row = row.Next()
	}
	sort.Sort(sortByTaskDeadline(tasks))
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
			link, exists := sel.Attr("href")
			if exists && !strings.Contains(link, "mailto") && !strings.Contains(link, "handins") {
				if strings.Contains(link, "given") {
					files = append(files, getGivenFiles(link)...)
				} else {
					files = append(files, link)
				}
			}
		},
	)
	//Remove redundant whitespace
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(sel.Text(), " ")
	return &Task{
		Name:      strings.TrimSpace(s),
		Class:     coloursToGroups[colour],
		Deadline:  convertDaysToDate(day),
		Links:     files,
		FileNames: make([]string, len(files)),
	}
}

//stores the module struct in modules.json
func storeModules() error {
	data, err := json.MarshalIndent(Modules, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path+"/"+modulePath, data, 0644)
}

func loadModules() error {
	data, err := ioutil.ReadFile(path + "/" + modulePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Modules)
}

type sortByTaskDeadline []*Task

func (a sortByTaskDeadline) Len() int           { return len(a) }
func (a sortByTaskDeadline) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByTaskDeadline) Less(i, j int) bool { return a[i].Deadline.Before(a[j].Deadline) }
