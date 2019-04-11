package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Reaction hush
type Reaction struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	Count int      `json:"count"`
}

// Message hush
type Message struct {
	Reactions []Reaction `json:"reactions"`
}

// ReactionsGetResponse hush
type ReactionsGetResponse struct {
	OK      bool     `json:"ok"`
	Channel string   `json:"channel"`
	Message *Message `json:"message"`
}

// Channel hush
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChannelListResponse hush
type ChannelListResponse struct {
	Channels []Channel `json:"channels"`
}

func fetchChannelList() {
	slackAPIUrl := "https://slack.com/api/"
	channelsList := "channels.list"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(channelsList)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	q.Add("token", "INSERT TOKEN HERE")
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.

	if err != nil {
		log.Fatalf("Failed to create request %v", err)
	}

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatalf("Failed to fetch channels %v", getErr)
	}

	var slackMsg ChannelListResponse
	err = json.NewDecoder(res.Body).Decode(&slackMsg)
	if err != nil {
		log.Fatalf("Failed to decode %v", err)
	}
	fmt.Println("Got these channels back.")
	fmt.Println(slackMsg.Channels)
}

// conversationsList := "channels.list"

func fetchEmojis() {
	slackAPIUrl := "https://slack.com/api/"
	reactionsGet := "reactions.get"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(reactionsGet)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	q.Add("token", "xoxp-2322548031-351754944855-605954585685-8730f3f5dc586dd31575c841b52364eb")
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

	var slackMsg ReactionsGetResponse
	err = json.NewDecoder(res.Body).Decode(&slackMsg)
	if err != nil {
		log.Fatalf("Failed to decode %v", err)
	}
	fmt.Println("Got these reactions back.")
	fmt.Println(slackMsg.Message.Reactions)
}

func main() {
	// fetchEmojis()
	fetchChannelList()
}
