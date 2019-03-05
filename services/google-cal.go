package calevents

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
// channel that recieves the map of event names, id
var idMapCh = make(chan map[string]string)
var mapReady = make(chan bool)

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
	_, ok := eventMap[strings.ToLower(eventName)]
	// if NOT in map
	if !ok {
		eventMap[strings.ToLower(eventName)] = eventID // add to map name, id
	}
	return // break the function

}

//MakeIDMap makes map of event name, id pairs
// is a gouroutine that runs in the background every
// n minutes
func MakeIDMap() {
	var nameIDMap = make(map[string]string)
	srv := getService() // gets authorized cal service
	// grab current time to start look from
	t := time.Now().Format(time.RFC3339)

	go func() {
		// get list of events
		events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalln(err)
		}
		// loop over items and create map
		for _, item := range events.Items {
			// creates the map
			createEventMap(item.Summary, item.Id, nameIDMap)
		}
		// sends the nameIDMap throught the idMapCh
		idMapCh <- nameIDMap // sends new map to idMapCh
	}()

}

//UpdateMap is a goroutine that keeps the event, id map updated
// it reads from the idMapCh
func UpdateMap() {
	go func() {
		eMap := <-idMapCh
		eventMap = eMap // updates the global event map
		mapReady <- true
	}()
}

//Index returns slice of all upcoming event objects
func Index() ([][]byte, error) {
	var eventSlice [][]byte // empty slice for events
	// gets the authorized cal service
	srv := getService()
	// gets the current time to assign as when to start looking for events
	t := time.Now().Format(time.RFC3339)
	// gets list of events from google cal package
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		return nil, err
	}
	type ItemTime struct {
		DateTime time.Time `json:"dateTime"`
	}
	// var startObj ItemTime
	// var endObj ItemTime
	// check if no events
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			// s, _ := item.Start.MarshalJSON()
			// e, _ := item.End.MarshalJSON()
			// json.Unmarshal(e, &endObj)
			// json.Unmarshal(s, &startObj)

			itm, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			// appaned marshalled event item into the event slice
			eventSlice = append(eventSlice, itm)
		}
	}
	return eventSlice, nil
}

//Add takes in an event info and adds to google cal
func Add(summary string, description string, start string, end string, colorID string) {
	// srv := getService() // gets google service

	//event struct to be used for calendar insert call
	type eventObj struct {
		summary     string    // title of event
		description string    // description
		start       time.Time // start time of event
		end         time.Time // end time of event
		colorID     string    // color of event card
	}

}

//QuickAdd takes a string and adds an hour event at current time
func QuickAdd(eventText string) {
	srv := getService() // gets google cal service

	newEvent := srv.Events.QuickAdd("primary", eventText)
	_, err1 := newEvent.Do()
	if err1 != nil {
		panic(err1)
	}

	fmt.Printf("New Event Created: %s", eventText)
}

//Remove removes specified event from your calendar
func Remove(eventName string) error {
	<-mapReady          // blocks from being called until idMap is created
	srv := getService() // get google cal service

	// find id using eventMap and inputted name
	id, ok := eventMap[strings.ToLower(eventName)]
	if !ok {
		fmt.Printf("%s not in calendar", eventName)
		err := srv.Events.Delete("primary", id).Do()
		if err != nil {
			return err
		}
		fmt.Printf("Event: %s, was successfully deleted!", eventName)
		return nil // no error return nothing
	}
	fmt.Println("event not in calendar")
	return nil
}

//Find returns information about specific event (if exists)
func Find(eventName string) ([]byte, error) {
	<-mapReady          // blocks this from being called until idMap is created
	srv := getService() // gets google cal service
	// use event map and get id
	id, ok := eventMap[strings.ToLower(eventName)]
	// check to see if event was in map
	if !ok {
		fmt.Printf("EVENT: %s, not found in calendar.", eventName)
		return nil, nil //! returning out with err I assume
	}
	// get the event front google cal
	eventObj, err := srv.Events.Get("primary", id).Do()
	if err != nil {
		return nil, err
	}
	event, err := json.Marshal(eventObj)
	if err != nil {
		return nil, err
	}

	// return the event json object
	return event, nil
}
