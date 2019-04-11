// Following this tutorial: https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
// Uses this router: https://github.com/gorilla/mux
// Run with: go build && ./rest-api

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Person struct
type Person struct {
	// Struct tags: https://godoc.org/encoding/json#Marshal
	// https://github.com/golang/go/wiki/Well-known-struct-tags
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firstname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// Address - a Person's address
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

// GetPeople returns all known Person entities
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetPerson finds an existing Person given their ID
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, person := range people {
		if person.ID == params["id"] {
			json.NewEncoder(w).Encode(person)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// CreatePerson creates a new peson entity
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Printf("Error decoding person: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson deletes an existing person entity
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			// https://programming.guide/go/three-dots-ellipsis.html
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}

type test struct {
	Number int `json:"number"`
}

// Reaction { "name": "tada", "users": [ "UABN6TSR5" ], "count": 1 },
type Reaction struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	Count int      `json:"count"`
}

// Message hush
type Message struct {
	Reactions []Reaction `json:"reactions"`
}

// SlackAPIResponse it is what it is.
type SlackAPIResponse struct {
	OK      bool     `json:"ok"`
	Channel string   `json:"channel"`
	Message *Message `json:"message"`
}

func fetchEmojis() {
	slackAPIUrl := "https://slack.com/api/"
	reactionsPath := "reactions.get"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(reactionsPath)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	q.Add("token", "xoxp-2322548031-351754944855-609816995558-5ff5c130f4f57291c3ebe25da55751e9")
	q.Add("channel", "CHGG9F6BU")
	q.Add("timestamp", "1555014521.000700")
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.

	if err != nil {
		log.Fatalf("Failed to create request %v", err)
	}

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatalf("Failed to fetch emojis %v", getErr)
	}

	var slackMsg SlackAPIResponse
	err = json.NewDecoder(res.Body).Decode(&slackMsg)
	if err != nil {
		log.Fatalf("Failed to decode %v", err)
	}
	fmt.Println("Got these reactions back.")
	fmt.Println(slackMsg.Message.Reactions)
}

func main() {
	fetchEmojis()
}
