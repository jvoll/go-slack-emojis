// Following this tutorial: https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
// Uses this router: https://github.com/gorilla/mux
// Run with: go build && ./rest-api

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type SlackApiResponse struct {
	OK      bool   `json:"ok"`
	Channel string `json:"channel"`
	// Message *Message `json:"message"`
	// Reactions []Reaction `json:"reactions"`
}

func FetchEmojis() {
	// (originally) Stolen from https://blog.alexellis.io/golang-json-api-client/
	// url := "http://api.open-notify.org/astros.json"
	url := "https://slack.com/api/reactions.get?token=xoxp-2322548031-380097601427-605976644912-eaf3c87c6e9ce58cc3615f46c7daffc8&channel=CHGG9F6BU&timestamp=1555014521.000700&pretty=1"
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Failed to create request %v", err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatalf("Failed to fetch emojis %v", getErr)
	}

	var slackMsg SlackApiResponse
	err = json.NewDecoder(res.Body).Decode(&slackMsg)
	if err != nil {
		log.Fatalf("Failed to decode %v", err)
	}
	fmt.Println("Got this channel id back")
	fmt.Println(slackMsg.Channel)
}

func main() {
	FetchEmojis()

	people = append(people, Person{ID: "1", FirstName: "Jason", LastName: "Voll", Address: &Address{City: "Cambridge", State: "MA"}})
	people = append(people, Person{ID: "2", FirstName: "Payal", LastName: "Batra", Address: &Address{City: "Cambridge", State: "MA"}})
	people = append(people, Person{ID: "3", FirstName: "Dan", LastName: "Dexter"})

	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	port := ":8000"
	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, router))
}
