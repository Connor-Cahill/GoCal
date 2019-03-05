package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	calevents "github.com/connor-cahill/goCal/services"
	"github.com/spf13/cobra"
)

//Event struct for displaying events
type Event struct {
	Summary string    `json:"summary"`
	Start   EventTime `json:"start"`
	End     EventTime `json:"end"`
}

//EventTime is a time object for the event
type EventTime struct {
	Time string `json:"dateTime"`
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
			err = json.Unmarshal(e, &event)
			if err != nil {
				log.Fatalln(err)
			}
			if event.Start.Time == "" {
				fmt.Println("BREAKING:  ", event.Summary)

			} else {
				st := strings.Split(event.Start.Time, "-")
				et := strings.Split(event.End.Time, "-")
				_ = i
				st = st[:len(st)-1]
				et = et[:len(et)-1]
				fmt.Println(event.Start.Time)
				fmt.Println(strings.Join(st, "-"))
				sTime, err := time.Parse("2006-01-02T15:04:05", strings.Join(st, "-"))
				if err != nil {
					log.Fatalln(err)
				}
				eTime, err := time.Parse("2006-01-02T15:04:05", strings.Join(et, "-"))
				if err != nil {
					log.Fatalln(err)
				}
				var startTime []int
				var endTime []int
				hr, min, _ := sTime.Clock()
				startTime = append(startTime, hr)
				startTime = append(startTime, min)
				hr, min, _ = eTime.Clock()
				endTime = append(endTime, hr)
				endTime = append(endTime, min)
				fmt.Println(startTime)
				fmt.Println(endTime)
				fmt.Println("===========================")
				// prints 1. EventSummary
				// fmt.Printf("%d. %s\n STARTS AT: %s - ENDS AT: %s\n", i+1, event.Summary, start, end)
			}
		}
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(showCmd)
}
