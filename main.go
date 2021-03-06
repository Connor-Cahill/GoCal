package main

import (
	"log"

	"github.com/connor-cahill/goCal/cmd"
	calevents "github.com/connor-cahill/goCal/services"
)

func main() {

	// creates the event name, id map
	calevents.MakeIDMap()
	// creates event map with event names, ids
	calevents.UpdateMap()
    // Auths the user if not already
    err := calevents.AuthUser()
    if err != nil {
        log.Fatal(err)
    }
	//* sets up CLI with cobra
	err = cmd.RootCmd.Execute() // execute our CLI package
	if err != nil {
		log.Fatalln(err)
	}
    
}
