package main

import (
	"fmt"
	"time"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/Akshat-Tripathi/cateCli/fileopen"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

//List prints out all the modules or all the tasks of a module
//If showTask is true then module shouldn't be an empty string
func List(showTask bool, module string) {
	if !showTask {
		listModules()
	} else {
		listTasks(module)
	}
}

func listModules() {
	for _, mod := range cate.Modules {
		fmt.Println(mod.Name)
	}
}

func listTasks(module string) {
	mod := getModule(module)
	time.Sleep(time.Millisecond * 50) //Without this there's a chance that some text doesn't print
	for _, task := range mod.Tasks {
		fmt.Println(colourTaskName(task))
	}
}

func colourTaskName(task *cate.Task) string {
	if task.Deadline.Before(time.Now()) {
		return color.BlueString(task.Name)
	}
	switch task.Class {
	case cate.Group:
		return color.HiMagentaString(task.Name)
	case cate.Individual:
		return color.HiGreenString(task.Name)
	case cate.Unassessed:
		return color.HiWhiteString(task.Name)
	case cate.UnassessedSub:
		return color.HiCyanString(task.Name)
	default:
		return color.BlackString("")
	}

}

func getModule(mod string) *cate.Module {
	if mod == "" {
		mod = selectModule()
	}
	module, err := findModule(mod)
	if err != nil {
		panic("Module not found, try running fetch")
	}
	return module
}

func selectModule() string {
	modules := make([]string, len(cate.Modules))
	for i, mod := range cate.Modules {
		modules[i] = mod.Name
	}
	prompt := promptui.Select{
		Label: "Select a module",
		Items: modules,
		Size:  len(modules),
	}

	_, module, err := prompt.Run()
	if err != nil {
		panic("Couldn't select module")
	}
	return module
}

func getTask(task string, mod *cate.Module) *cate.Task {
	if task == "" {
		task = selectTask(mod)
	}
	for _, tsk := range mod.Tasks {
		if tsk.Name == task {
			return tsk
		}
	}
	panic("Task not found")
}

func selectTask(mod *cate.Module) string {
	tasks := make([]string, 0)
	for _, task := range mod.Tasks {
		if len(task.Links) > 0 {
			tasks = append(tasks, task.Name)
		}
	}
	prompt := promptui.Select{
		Label: "Select a Task",
		Items: tasks,
	}

	_, task, err := prompt.Run()
	if err != nil {
		panic("Couldn't select task")
	}
	return task
}

//Get downloads all the files needed for a particular task
func Get(module, task string) {
	mod := getModule(module)
	tsk := getTask(task, mod)
	cate.DownloadTask(tsk, mod)
}

//Show opens all the files in a task
func Show(module, task string) {
	mod := getModule(module)
	tsk := getTask(task, mod)
	loc := cate.ModulePath(mod)
	for _, name := range tsk.FileNames {
		fileopen.Open(loc, name)
	}
}
