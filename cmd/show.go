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
                fmt.Printf("%d. %s -- All Day Event", i + 1, event.Summary)
			} else {
				st := strings.Split(event.Start.Time, "-")
				et := strings.Split(event.End.Time, "-")
				_ = i
				st = st[:len(st)-1]
				et = et[:len(et)-1]
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
                if hr > 12 {
                    hr -= 12
                }
				startTime = append(startTime, hr)
				startTime = append(startTime, min)
				hr, min, _ = eTime.Clock()
				if hr > 12 {
                    hr -= 12
                }
				endTime = append(endTime, hr)
				endTime = append(endTime, min)
                starting := strings.Trim(strings.Join(strings.Split(fmt.Sprint(startTime), " "), ":"), "[]")
                ending := strings.Trim(strings.Join(strings.Split(fmt.Sprint(endTime), " "), ":"), "[]")
				// prints 1. EventSummary
				fmt.Printf("%d. %s -- %s - %s\n", i+1, event.Summary, starting, ending)
			}
		}
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(showCmd)
}
