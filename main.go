package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var TOKEN = "TOKEN_HERE"

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

func fetchChannelList() []Channel {
	slackAPIUrl := "https://slack.com/api/"
	channelsList := "channels.list"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(channelsList)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	q.Add("token", TOKEN)
	q.Add("limit", "10")          // Limit us to 10 channels to not break the internet
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

	return slackMsg.Channels
}

// ChannelHistoryResponse hush
type ChannelHistoryResponse struct {
	Messages []Message `json:"messages"`
}

func fetchChannelHistory(channelID string) []Message {
	slackAPIUrl := "https://slack.com/api/"
	channelsHistory := "channels.history"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(channelsHistory)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	q.Add("token", TOKEN)
	q.Add("channel", channelID)
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

	var channelHistory ChannelHistoryResponse
	err = json.NewDecoder(res.Body).Decode(&channelHistory)
	if err != nil {
		log.Fatalf("Failed to decode %v", err)
	}
	return channelHistory.Messages
}

// func fetchEmojis() {
// 	slackAPIUrl := "https://slack.com/api/"
// 	reactionsGet := "reactions.get"

// 	var sb strings.Builder
// 	sb.WriteString(slackAPIUrl)
// 	sb.WriteString(reactionsGet)
// 	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
// 	q := req.URL.Query() // Get a copy of the query values.
// 	q.Add("token", TOKEN)
// 	q.Add("channel", "CHGG9F6BU")
// 	q.Add("timestamp", "1555014521.000700")
// 	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.

// 	if err != nil {
// 		log.Fatalf("Failed to create request %v", err)
// 	}

// 	client := http.Client{
// 		Timeout: time.Second * 2, // Maximum of 2 secs
// 	}
// 	res, getErr := client.Do(req)
// 	if getErr != nil {
// 		log.Fatalf("Failed to fetch emojis %v", getErr)
// 	}

// 	var slackMsg ReactionsGetResponse
// 	err = json.NewDecoder(res.Body).Decode(&slackMsg)
// 	if err != nil {
// 		log.Fatalf("Failed to decode %v", err)
// 	}
// 	fmt.Println("Got these reactions back.")
// 	fmt.Println(slackMsg.Message.Reactions)
// }

func main() {
	// fetchEmojis()

	reactions := make(map[string]int)

	channels := fetchChannelList()
	for _, channel := range channels {
		recentMessages := fetchChannelHistory(channel.ID)

		for _, message := range recentMessages {
			// Only look at messages with emoji reactions
			if message.Reactions != nil {

				for _, reaction := range message.Reactions {
					_, hasKey := reactions[reaction.Name]
					if hasKey {
						reactions[reaction.Name] += reaction.Count
					} else {
						reactions[reaction.Name] = reaction.Count
					}
				}
			}
		}
	}

	fmt.Println(reactions)
}
