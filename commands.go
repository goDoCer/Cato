package main

import (
	"fmt"
	"time"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

//List prints out all the modules or all the tasks of a module
//If showTask is true then module shouldn't be an empty string
func list(showTask bool, module string) {
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
	if module == "" {
		module = selectModule()
	}
	mod, err := findModule(module)
	if err != nil {
		panic("Module not found, try running fetch")
	}
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

func selectModule() string {
	modules := make([]string, len(cate.Modules))
	for i, mod := range cate.Modules {
		modules[i] = mod.Name
	}
	prompt := promptui.Select{
		Label: "Select a module",
		Items: modules,
	}

	_, module, err := prompt.Run()
	if err != nil {
		panic("Couldn't select module")
	}
	return module
}

func selectTask(mod *cate.Module) string {
	tasks := make([]string, len(mod.Tasks))
	for i, task := range mod.Tasks {
		tasks[i] = task.Name
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

func get(module, task string) {

}
