package main

import (
	"os"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/urfave/cli"
)

func main() {
	cate.Init()

	app := &cli.App{
		Commands: []*cli.Command{
			Fetch(),
			Ls(),
			Login(),
			Get(),
			Show(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	// switch os.Args[1] {
	// case "init":
	// 	cate.Init()
	// case "fetch":
	// 	cate.Fetch()
	// case "ls":
	// 	List(len(os.Args) == 3, "")
	// case "login":
	// 	cate.Login()
	// case "get":
	// 	Get(getArg(2, os.Args), getArg(3, os.Args))
	// case "show":
	// 	Show(getArg(2, os.Args), getArg(3, os.Args))
	// default:
	// 	os.Exit(1)
	// }
}
