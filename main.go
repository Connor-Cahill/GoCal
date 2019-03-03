package main

import (
	"log"

	"github.com/connor-cahill/goCal/cmd"
	"github.com/connor-cahill/goCal/services"
)

func main() {

	// creates event map with event names, ids
	startList, err := services.Index()
	if err != nil {
		log.Fatalln(err)
	}
	_ = startList // creates the initial event map

	//* sets up CLI with cobra
	err = cmd.RootCmd.Execute() // execute our CLI package
	if err != nil {
		log.Fatalln(err)
	}

}
