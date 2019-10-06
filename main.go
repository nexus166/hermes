package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	rocketChatURL     string
	rocketChatToken   string
	rocketChatTokenId string
	rocketChatChannel string
	rocketChatRoomId  string
)

var post struct {
	Text  string
	Emoji string
}

type postTpl struct {
	Channel string `json:"channel,omitempty"`
	RoomId  string `json:"roomId,omitempty"`
	Text    string `json:"text"`
	Emoji   string `json:"emoji,omitempty"`
}

func init() {

	if rocketChatTokenEnv, set := os.LookupEnv("ROCKETCHAT_TOKEN"); set {
		fmt.Println("Using env variable ROCKETCHAT_TOKEN")
		rocketChatToken = rocketChatTokenEnv
	}
	if rocketChatToken == "" {
		fmt.Println("Need a token to run..")
		os.Exit(127)
	}

	if rocketChatTokenIdEnv, set := os.LookupEnv("ROCKETCHAT_TOKEN_ID"); set {
		fmt.Println("Using env variable ROCKETCHAT_TOKEN_ID")
		rocketChatTokenId = rocketChatTokenIdEnv
	}
	if rocketChatTokenId == "" {
		fmt.Println("Need a token ID to run..")
		os.Exit(127)
	}

	if rocketChatURLEnv, set := os.LookupEnv("ROCKETCHAT_URL"); set {
		fmt.Println("Using env variable ROCKETCHAT_URL")
		rocketChatURL = rocketChatURLEnv
	}
	if rocketChatURL == "" {
		fmt.Println("Need a target RocketChat instance..")
		os.Exit(127)
	}
	rocketChatURL = rocketChatURL + "/api/v1/chat.postMessage"

	channelFlag := flag.String("c", "", "channel to post message in")
	roomIdFlag := flag.String("r", "", "room to post message in")

	flag.StringVar(&post.Text, "m", "test message", "message to post")
	flag.StringVar(&post.Emoji, "e", "", "emoji to attach")
	flag.Parse()

	if *channelFlag != "" {
		rocketChatChannel = *channelFlag
	}
	if *roomIdFlag != "" {
		rocketChatRoomId = *roomIdFlag
	}
	if rocketChatChannel == "" && rocketChatRoomId == "" {
		fmt.Println("Need a RoomId or Channel to post the message in..")
		os.Exit(127)
	}

	if post.Text == "-" || post.Text == "" {
		if console, _ := os.Stdin.Stat(); (console.Mode() & os.ModeCharDevice) == 0 {
			stdinmsg, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				return
			}
			post.Text = string(stdinmsg)
		}
	}
	if post.Text == "" {
		fmt.Println("No data in message, sending test message..")
		post.Text = "It works!"
	}
}

func main() {
	payload := &postTpl{
		Channel: rocketChatChannel,
		RoomId:  rocketChatRoomId,
		Text:    post.Text,
		Emoji:   post.Emoji,
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(payload)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b.Bytes()))

	req, err := http.NewRequest("POST", rocketChatURL, b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", rocketChatToken)
	req.Header.Set("X-User-Id", rocketChatTokenId)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
