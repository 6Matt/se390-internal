package main

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
)

// Helpers
const endpoint = "http://clashfinder.com/data/"
const ctLayout = "2006-01-02 15:04"

func getJson(url string, target interface{}) error {
	fmt.Println("getting json from:", url)
	r, err := http.Get(url)
	if err != nil {
    	return err
	}
	defer r.Body.Close()
	// fmt.Println("response:", r)
	return json.NewDecoder(r.Body).Decode(target)
}

type CustomTime struct {
    time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
    if b[0] == '"' && b[len(b)-1] == '"' {
        b = b[1 : len(b)-1]
    }
    ct.Time, err = time.Parse(ctLayout, string(b))
    return
}

type Festival struct {
    Id 		string	`json:"name"`
    Name 	string	`json:"desc"`
    IsCore	bool	`json:"coreClashfinder"`
}

type Event struct {
	Name 	string
	Start 	time.Time
	End 	time.Time
}

type Location struct {
	Name 	string
	Events 	[]Event
}

// Get all Festivals
func getFestivals() []Festival {
	m := map[string]Festival{}
	url := endpoint + "events/all.json"
    error := getJson(url, &m)
	if error != nil {
		fmt.Println(error)
	}
	festivals := make([]Festival, 0, len(m))
	for _, v := range m {
    	festivals = append(festivals, v)
    }
	return festivals
}

// Get an Event schedule
func getSchedule(id string) []Location {
	type EventCustTime struct {
		Name 	string
		Start 	CustomTime
		End 	CustomTime
	}

	type Response struct {
		Locations []Location
	}

	schedule := Response{}
	url := endpoint + "event/" + id + ".json"
    error := getJson(url, &schedule)
	if error != nil {
		fmt.Println(error)
	}
	return schedule.Locations
}

// Simple main for testing
func main() {
	festivals := getFestivals()
	fmt.Println("\nFESTIVALS:\n\n", festivals, "\n")
	sched := getSchedule("osheaga2016official")
	fmt.Println("\nSCHEDULE:\n\n", sched, "\n")
}
