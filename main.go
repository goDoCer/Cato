package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Akshat-Tripathi/cateCli/cate"
)

func main() {
	cate.Init()
	getCommand := flag.NewFlagSet("get", flag.ExitOnError)

	//get subcommands
	getPtr := getCommand.String("module", "", "Module to get (Required)")

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cate.Init()
	case "fetch":
		cate.Fetch()
	case "get":
		getCommand.Parse(os.Args[2:])
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
	}
}
