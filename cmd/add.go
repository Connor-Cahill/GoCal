package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new event to your calendar",
	Run: func(cmd *cobra.Command, args []string) {
		// desc, _ := getEventAndTime()
		// description := []byte(*desc)
		description := strings.Join(args, " ")
		fmt.Println(description)
		// services.AddNewEvent(description)
	},
}

// func getEventAndTime() (*string, *string) {
// 	description := flag.String("event description", "Please enter event description: ", "gets description of event to add to calendar")
// 	eTime := flag.String("event time", "Please enter a time for this event: ", "Gets the time for the event.")
// 	return description, eTime
// }

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(addCmd)
}
