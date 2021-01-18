package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Akshat-Tripathi/cateCli/cate"
)

func main() {
	cate.Init()
	getCommand := flag.NewFlagSet("get", flag.ExitOnError)
	showCommand := flag.NewFlagSet("show", flag.ExitOnError)

	//get subcommands
	getPtr := getCommand.String("module", "", "Module to get (Required)")
	showModulePtr := showCommand.String("module", "", "Module to show (Required)")

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cate.Init()
	case "fetch":
		cate.Fetch()
	case "ls":
		list(true, "")
	case "login":
		cate.Login()
	case "get":
		getCommand.Parse(os.Args[2:])
	case "show":
		showCommand.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}

	if getCommand.Parsed() {
		if *getPtr == "" {
			getCommand.PrintDefaults()
			os.Exit(1)
		}
		found := false
		for _, module := range cate.Modules {
			if *getPtr == strings.Split(module.Name, " ")[0] {
				cate.DownloadModule(module)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("No module matching:", *getPtr, "found")
		}
	} else if showCommand.Parsed() {
		if *showModulePtr == "" {
			getCommand.PrintDefaults()
			os.Exit(1)
		}
		module, err := findModule(*showModulePtr)
		if err != nil {
			panic(err)
		}
		filepath.Walk("files/"+module.Name, func(path string, _ os.FileInfo, _ error) error {
			if strings.HasSuffix(path, "pdf") {
				cmd := exec.Command("explorer.exe", path)
				cmd.Run()
				fmt.Println("file", path)
			}
			return nil
		})
	}
}

func findModule(name string) (module *cate.Module, err error) {
	for _, v := range cate.Modules {
		if strings.Split(v.Name, " ")[0] == name || v.Name == name {
			return v, nil
		}
	}
	return nil, errors.New("Couldn't find module " + name)
}
