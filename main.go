package main

import (
	"os"
	"strings"

	"github.com/Akshat-Tripathi/cateCli/cate"
	"github.com/kardianos/osext"
	"github.com/urfave/cli/v2"
)

var path string

func main() {
	p, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	path = strings.ReplaceAll(p, "\\", "/")
	cate.Init(path)

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

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
