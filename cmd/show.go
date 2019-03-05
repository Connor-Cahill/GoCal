package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	calevents "github.com/connor-cahill/goCal/services"
	"github.com/spf13/cobra"
)

//Event struct for displaying events
type Event struct {
	Summary string  `json:"summary"`
	//Start   string  `json:"start,omitempty"`
	//End     string  `json:"end,omitempty"`
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "prints all upcoming events to terminal",
	Run: func(cmd *cobra.Command, args []string) {
		events, err := calevents.Index()
		if err != nil {
			log.Fatalln(err)
		}
		// event to be unmarshaled
		var event Event
		// list of events to print to terminal
		// var eventList []Event
		for i, e := range events {
			err := json.Unmarshal(e, &event)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("%d. %s\n", i+1, event.Summary)

		}
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(showCmd)
}
