package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/Akshat-Tripathi/cateCli/fileopen"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

//Fetch downloads the current term's timetable
func Fetch() cli.Command {
	return cli.Command{
		Name:  "fetch",
		Usage: "gets the current timetable information",
		Action: func(c *cli.Context) error {
			return cate.Fetch()
		},
	}
}

//Get downloads all the files needed for a particular task
func Get() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "downloads all files related to a task",
		Action: func(c *cli.Context) error {
			module := c.Args().Get(2)
			task := c.Args().Get(3)
			mod, err := getModule(module)
			if err != nil {
				return err
			}
			tsk, err := getTask(task, mod)
			if err != nil {
				return err
			}
			return cate.DownloadTask(tsk, mod)
		},
	}
}

//Show opens all the files in a task
func Show() cli.Command {
	return cli.Command{
		Name:  "show",
		Usage: "opens all files related to a task",
		Action: func(c *cli.Context) error {
			module := c.Args().Get(2)
			task := c.Args().Get(3)
			mod, err := getModule(module)
			if err != nil {
				return err
			}
			tsk, err := getTask(task, mod)
			if err != nil {
				return err
			}
			loc := cate.ModulePath(mod)
			for _, name := range tsk.FileNames {
				err = fileopen.Open(loc, name)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

//Ls prints out all the modules or all the tasks of a module
//If showTask is true then module shouldn't be an empty string
func Ls() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "lists modules and/or tasks",
		Flags: []cli.Flag{
			&cli.BoolTFlag{
				Name:  "task",
				Usage: "shows tasks in a module",
			},
		},
		Action: func(c *cli.Context) error {
			module := c.Args().Get(2)
			if c.BoolT("task") {
				return listTasks(module)
			} else {
				listModules()
			}
			return nil
		},
	}
}

//Login saves the users login in secrets.json
func Login() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "save login details",
		Action: func(c *cli.Context) error {
			return cate.Login()
		},
	}
}

func listModules() {
	for _, mod := range cate.Modules {
		fmt.Println(mod.Name)
	}
}

func listTasks(module string) error {
	mod, err := getModule(module)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 50) //Without this there's a chance that some text doesn't print
	for _, task := range mod.Tasks {
		fmt.Println(colourTaskName(task))
	}
	return nil
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

func getModule(mod string) (*cate.Module, error) {
	if mod == "" {
		mod = selectModule()
	}
	module, err := findModule(mod)
	if err != nil {
		return nil, err
	}
	return module, nil
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

func getTask(task string, mod *cate.Module) (*cate.Task, error) {
	if task == "" {
		task = selectTask(mod)
	}
	for _, tsk := range mod.Tasks {
		if tsk.Name == task {
			return tsk, nil
		}
	}
	return nil, fmt.Errorf("Task - %s not found", task)
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

func findModule(name string) (module *cate.Module, err error) {
	for _, v := range cate.Modules {
		if strings.Split(v.Name, " ")[0] == name || v.Name == name {
			return v, nil
		}
	}
	return nil, errors.New("Couldn't find module " + name)
}
