package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nexus166/hermes"
)

func main() {
	if err := hermes.CLIFlags.Parse(os.Args[1:]); err != nil && err != flag.ErrHelp {
		panic(err)
	}
	if hermes.RCChannel == "" && hermes.RCRoomID == "" {
		panic("need a RoomID or Channel to post the message in..")
	}
	if hermes.RCToken == "" || hermes.RCTokenID == "" {
		panic("need a token/ID to run..")
	}
	if hermes.RCURL == "" {
		panic("need a target RC instance..")
	}
	data := strings.Join(hermes.CLIFlags.Args(), " ")
	if data == "-" {
		if console, _ := os.Stdin.Stat(); (console.Mode() & os.ModeCharDevice) == 0 {
			stdinmsg, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				return
			}
			data = string(stdinmsg)
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if data = scanner.Text(); data != "" {
					break
				}
			}
		}
	}
	if data == "" {
		fmt.Println("no data in message, sending test message..")
		data = "It works!"
	}
	if status, err := hermes.POSTMessage(hermes.POSTTemplate{
		Channel: hermes.RCChannel,
		RoomID:  hermes.RCRoomID,
		Text:    data,
		Emoji:   hermes.RCEmoji,
	}); err != nil {
		fmt.Println(status)
		panic(err)
	}
}
