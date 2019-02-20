package main

import (
	"github.com/connor-cahill/goCal/cmd"
)

func main() {
	//* sets up CLI with cobra
	err := cmd.RootCmd.Execute() // execute our CLI package
	if err != nil {
		panic(err)
	}
}
