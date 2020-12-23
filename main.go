package main

import (
	"github.com/Akshat-Tripathi/cateCli/parser"
)

func main() {
	// if err := parser.Download(); err != nil {
	// 	panic(err)
	// }
	// parser.Init()
	// if err := parser.DownloadTimeTable(); err != nil {
	// 	panic(err)
	// }
	doc := parser.L()
	parser.GetModules(doc)
}
