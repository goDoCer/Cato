package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/Akshat-Tripathi/cateCli/fileopen"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

//Fetch downloads the current term's timetable
func Fetch() *cli.Command {
	return &cli.Command{
		Name:  "fetch",
		Usage: "gets the current timetable information",
		Action: func(c *cli.Context) error {
			return cate.Fetch()
		},
	}
}

//Get downloads all the files needed for a particular task
func Get() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "downloads all files related to a task",
		Action: func(c *cli.Context) error {
			module := c.Args().Get(0)
			task := c.Args().Get(1)
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
func Show() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "opens all files related to a task",
		Action: func(c *cli.Context) error {
			module := c.Args().Get(0)
			task := c.Args().Get(1)
			mod, err := getModule(module)
			if err != nil {
				return err
			}
			tsk, err := getTask(task, mod, func(t *cate.Task) bool {
				return t.Downloaded
			})
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
func Ls() *cli.Command {
	return &cli.Command{
		Name:  "ls",
		Usage: "lists modules and/or tasks",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "task",
				Usage:   "shows tasks in a module",
				Aliases: []string{"t"},
			},
			&cli.BoolFlag{
				Name:    "deadline",
				Usage:   "displays task deadlines",
				Aliases: []string{"d"},
			},
		},
		Action: func(c *cli.Context) error {
			if len(cate.Modules) == 0 {
				return fmt.Errorf("No modules found")
			}
			module := c.Args().Get(0)
			if c.Bool("task") {
				return listTasks(module, c.Bool("deadline"))
			}
			listModules()
			return nil
		},
	}
}

//Login saves the users login in secrets.json
func Login() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "save login details",
		Action: func(c *cli.Context) error {
			cate.GetLoginDetails()
			return cate.Login()
		},
	}
}

func listModules() {
	for _, mod := range cate.Modules {
		fmt.Println(mod.Name)
	}
}

func listTasks(module string, showDeadline bool) error {
	mod, err := getModule(module)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 50) //Without this there's a chance that some text doesn't print
	for _, task := range mod.Tasks {
		fmt.Print(colourTaskName(task))
		if showDeadline {
			fmt.Printf(" - %s\n", task.Deadline.Format(time.ANSIC))
		} else {
			fmt.Println()
		}
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
		return selectModule(cate.Modules), nil
	}
	modules, err := findModules(mod)
	if err != nil {
		return nil, err
	}
	if len(modules) == 1 {
		return modules[0], nil
	}
	return selectModule(modules), nil
}

func selectModule(mods []*cate.Module) *cate.Module {
	modules := make([]string, len(mods))
	for i, mod := range mods {
		modules[i] = mod.Name
	}
	prompt := promptui.Select{
		Label: "Select a module",
		Items: modules,
		Size:  len(modules),
	}

	moduleIndex, _, err := prompt.Run()
	if err != nil {
		panic("Couldn't select module")
	}
	return mods[moduleIndex]
}

func getTask(task string, mod *cate.Module, filters ...func(*cate.Task) bool) (*cate.Task, error) {
	var err error
	if task == "" {
		task, err = selectTask(mod, filters...)
	}
	if err != nil {
		return nil, err
	}
	for _, tsk := range mod.Tasks {
		if tsk.Name == task {
			return tsk, nil
		}
	}
	return nil, fmt.Errorf("Task - %s not found", task)
}

func selectTask(mod *cate.Module, filters ...func(*cate.Task) bool) (string, error) {
	tasks := make([]string, 0)
	var ok bool
	for _, task := range mod.Tasks {
		ok = true
		for _, filter := range filters {
			ok = ok && filter(task)
		}
		if len(task.Links) > 0 && ok {
			tasks = append(tasks, task.Name)
		}
	}
	if len(tasks) == 0 {
		return "", fmt.Errorf("No tasks found for module: %s", mod.Name)
	} else if len(tasks) == 1 {
		return tasks[0], nil
	}
	prompt := promptui.Select{
		Label: "Select a Task",
		Items: tasks,
	}

	_, task, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return task, nil
}

func findModules(name string) (modules []*cate.Module, err error) {
	matcher := regexp.MustCompile(fmt.Sprintf("(?i).*%s.*", name))
	modules = make([]*cate.Module, 0, len(cate.Modules))
	for _, v := range cate.Modules {
		if matcher.Match([]byte(v.Name)) {
			modules = append(modules, v)
		}
	}
	if len(modules) == 0 {
		return nil, errors.New("Couldn't find module " + name)
	}
	return modules, nil
}
