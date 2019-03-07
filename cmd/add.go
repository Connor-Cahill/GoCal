package cmd

import (
	"fmt"
	"bufio"
	"os"
    
    calevents "github.com/connor-cahill/goCal/services"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new event to your calendar",
	Run: func(cmd *cobra.Command, args []string) {
	    // create reader to read from standard input
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("Enter Event Title:  ")
        // saves the event title from terminal input
        eventTitle, _ := reader.ReadString('\n')
        fmt.Println("Enter a description of event:  ")
        // grabs the event description from terminal
        desc, _ := reader.ReadString('\n')
        fmt.Println("When will it start? (input ex. 12-03-2019 4:00pm ):  ")
        startingTime, _ := reader.ReadString('\n')
        fmt.Println("How long is the event(in hours)? (input ex. 2)") 
        endingTime, _ := reader.ReadString('\n')
        fmt.Println("=================================")
        fmt.Println(calevents.FixTime(startingTime)) 
        _ = eventTitle
        _ = desc
        _ = endingTime
	},
}


//	Adds show command to root command
func init() {
	RootCmd.AddCommand(addCmd)
}
