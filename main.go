package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

var slackAPIUrl = "https://slack.com/api/"

// Reaction - Contains the name of the reaction, the users who reacted with it,
//            and the number of times it was reacted with.
type Reaction struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	Count int      `json:"count"`
}

// Message - Contains a list of Reactions.
type Message struct {
	Reactions []Reaction `json:"reactions"`
}

// Channel - Contains the ID & Name of a channel.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChannelListResponse - Contains list of Channels.
type ChannelListResponse struct {
	Channels []Channel `json:"channels"`
}

// ChannelHistoryResponse - Contains an array of the last 100 messages from the channel.
type ChannelHistoryResponse struct {
	Messages []Message `json:"messages"`
}

func getToken() string {
	token, err := ioutil.ReadFile("oauth-access-token.txt")

	if err != nil {
		panic(err)
	}

	return string(token)
}

func fetchChannelList() []Channel {
	channelsList := "channels.list"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(channelsList)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	token := getToken()
	q.Add("token", token)
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

func fetchChannelHistory(channelID string) []Message {
	slackAPIUrl := "https://slack.com/api/"
	channelsHistory := "channels.history"

	var sb strings.Builder
	sb.WriteString(slackAPIUrl)
	sb.WriteString(channelsHistory)
	req, err := http.NewRequest(http.MethodGet, sb.String(), nil)
	q := req.URL.Query() // Get a copy of the query values.
	token := getToken()
	q.Add("token", token)
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

func main() {
	reactions := make(map[string]int)
	channels := fetchChannelList()
	// var m map[int]string
	// var keys []int
	// for k := range m {
	// 	keys = append(keys, k)
	// }
	// sort.Ints(keys)
	// for _, k := range keys {
	// 	fmt.Println("Key:", k, "Value:", m[k])
	// }

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

	fmt.Println("🐠", reactions)
	inverse := make(map[int]string)
	// numKeys := make([]int, len(reactions))
	for k, v := range reactions {
		inverse[v] += k + " | "
		// append(numKeys, v)
	}
	fmt.Println("🎈", inverse)
	// sort.Ints(numKeys)
	// keys := sort.Reverse(sort.IntSlice(numKeys))
	// fmt.Println("💎", numKeys)
	// fmt.Println("😱", keys)
	// for _, j := range numKeys {
	// 	fmt.Println(inverse[j])
	// }
}
