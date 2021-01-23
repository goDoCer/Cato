package main

import (
	"os"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/urfave/cli/v2"
)

func main() {
	cate.Init()

	app := &cli.App{
		Name: "Cato",
		Commands: []*cli.Command{
			Fetch(),
			Ls(),
			Login(),
			Get(),
			Show(),
		},
		UseShortOptionHandling: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
