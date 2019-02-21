# 🗓 GoCal - Google Calendar CLI 🗓
GoCal is a CLI that allows you to view, add, delete, update and find events on your Google Calendar. 

### Usage
All commands used from this CLI must be prefixed with gocal, example:

      $ gocal show

The above example is using GoCal's "show" command which returns a numbered list of all events on you calendar. Notice the "gocal" prefix. 

## Commands:
### Show
Returns a numbered list of all calendar events, example:

      $ gocal show

### Add
Adds a new event to your calendar at the current time, example:

      $ gocal add go to the gym

Notice everything after the "add" keyword is what the event will be named on your calendar.


### Tech Used:
- Go
- Cobra
- Google Calendar API

