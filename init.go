package hermes

import (
	"flag"
	"os"
)

// some/all these can be built into the binaries using ldflags
var (
	defaultRCURL     string
	defaultRCToken   string
	defaultRCTokenID string
	defaultRCChannel string
	defaultRCRoomID  string
	defaultRCEmoji   string = ":speaker:"
)

// get config from env
var (
	RCURL          = os.Getenv("ROCKETCHAT_URL")
	RCToken        = os.Getenv("ROCKETCHAT_TOKEN")
	RCTokenID      = os.Getenv("ROCKETCHAT_TOKEN_ID")
	RCChannel      = os.Getenv("ROCKETCHAT_CHANNEL")
	RCRoomID       = os.Getenv("ROCKETCHAT_ROOM_ID")
	RCEmoji        = os.Getenv("ROCKETCHAT_EMOJI")
	Verbose   bool = false
	CLIFlags  flag.FlagSet
)

func init() {
	if RCURL == "" && defaultRCURL != "" {
		RCURL = defaultRCURL
	}
	if RCToken == "" && defaultRCToken != "" {
		RCToken = defaultRCToken
	}
	if RCTokenID == "" && defaultRCTokenID != "" {
		RCTokenID = defaultRCTokenID
	}

	if RCChannel == "" && defaultRCChannel != "" {
		RCChannel = defaultRCChannel
	} else if u := os.Getenv("USER"); u != "" {
		RCChannel = "@" + u
	}
	CLIFlags.StringVar(&RCChannel, "c", RCChannel, "channel to post message in")

	if RCRoomID == "" && defaultRCRoomID != "" {
		RCRoomID = defaultRCRoomID
	}
	CLIFlags.StringVar(&RCRoomID, "r", RCRoomID, "room to post message in")

	if RCEmoji == "" && defaultRCEmoji != "" {
		RCEmoji = defaultRCEmoji
	}
	CLIFlags.StringVar(&RCEmoji, "e", RCEmoji, "emoji to attach")

	CLIFlags.BoolVar(&Verbose, "v", Verbose, "show HTTP request/response data")
}
