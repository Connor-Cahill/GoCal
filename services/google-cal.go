package services

//* NOTE: this file contains all setup and functions pertaining too google calendar api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

//* Map for users event names, id (key, value pair)
var eventMap = make(map[string]string) // event name, event Id

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

//timeIn converts time to local time of user
func timeIn(t time.Time, tzName string) time.Time {
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		panic(err)
	}
	return t.In(loc)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// returns a google service
func getService() *calendar.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalln(err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv // returns the authorized analytics service
}

//	creates a map of event name, event id key value pairs
func createEventMap(eventName string, eventID string, eventMap map[string]string) {

	// check if event name already exists in map
	_, ok := eventMap[eventName]
	// if NOT in map
	if !ok {
		eventMap[eventName] = eventID // add to map name, id
		return
	}
	return // already in map, do nothing

}

//GetEventList prints list of events as function
func GetEventList() ([]string, error) {
	var eventSlice []string // empty slice for events
	srv := getService()
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		return nil, err
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			itm, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			// create event map for summary, id pair
			createEventMap(item.Summary, item.Id, eventMap)
			eventSlice = append(eventSlice, string(itm))

		}
	}
	return eventSlice, nil
}

//AddEvent takes in an event info and adds to google cal
func AddEvent(summary string, description string, start string, end string, colorId string) {
	// srv := getService() // gets google service

	//event struct to be used for calendar insert call
	type eventObj struct {
		summary     string    // title of event
		description string    // description
		start       time.Time // start time of event
		end         time.Time // end time of event
		colorID     string    // color of event card
	}

	// TODO:
	// 1. parse the start and end time to be added into event struct
	// 2. finish this event var
	// 3. send event through Events.Insert call

	// var event = eventObj{
	// 	summary:     summary,
	// 	description: description,
	// }

	// srv.Events.Insert("primary")
}

//QuickAddNewEvent takes in event text and adds it to calendar
func QuickAddNewEvent(eventText string) {
	srv := getService() // gets google cal service

	newEvent := srv.Events.QuickAdd("primary", eventText)
	_, err1 := newEvent.Do()
	if err1 != nil {
		panic(err1)
	}

	fmt.Printf("New Event Created: %s", eventText)
}

//RemoveEvent removes specified event from your calendar
func RemoveEvent(eventName string) {
	srv := getService() // get google cal service

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	//	iterate through every event in events list
	for _, e := range events.Items {
		fmt.Println(e.Summary)
		if e.Summary == eventName { // if event summary equals inputted event name
			srv.Events.Delete("primary", e.Id).Do() // delete item from calendar
			fmt.Printf("Event: \"%s\" was successfully removed.", e.Summary)
			return // break the function
		}
	}
	fmt.Printf("EVENT: \"%s\" could not be found!", eventName) // no event was found matching inputted name
}

//FindSingleItem returns information about specific event (if exists)
func FindSingleItem(eventName string) {
	srv := getService() // gets google cal service

	t := time.Now().Format(time.RFC3339)
	// event, err := srv.Events.Get("primary", ) NOTE: this requires you to pass id.. how(?)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	for _, e := range events.Items {
		if strings.ToLower(e.Summary) == strings.ToLower(eventName) {
			desc := e.Description
			if desc == "" {
				desc = "{ There is no description for this event }"
			}
			fmt.Printf("Event Name  :  %s\nEvent Description  : %s\nStart Time  : %s\nEnd Time  : %s\n", e.Summary, desc, e.Start.DateTime, e.End.DateTime)
			return
		}
	}
	fmt.Printf("EVENT NOT FOUND: \"%s\" ", eventName)
}
