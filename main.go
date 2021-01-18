package main

import (
	"errors"
	"os"
	"strings"

	"github.com/Akshat-Tripathi/cateCli/cate"
)

func main() {
	cate.Init()
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cate.Init()
	case "fetch":
		cate.Fetch()
	case "ls":
		List(true, "")
	case "login":
		cate.Login()
	case "get":
		Get(getArg(2, os.Args), getArg(3, os.Args))
	default:
		os.Exit(1)
	}
}

func getArg(index int, args []string) string {
	if index+1 > len(args) {
		return ""
	}
	return args[index]
}

func findModule(name string) (module *cate.Module, err error) {
	for _, v := range cate.Modules {
		if strings.Split(v.Name, " ")[0] == name || v.Name == name {
			return v, nil
		}
	}
	return nil, errors.New("Couldn't find module " + name)
}
