package main

import (
	"os"
	"strings"
	"time"

	"github.com/jzelinskie/geddit"
	"github.com/nexus166/hermes"
	"github.com/nexus166/msg"
)

var (
	redditClientID     string  = os.Getenv("REDDIT_CLIENT_ID")
	redditClientSecret string  = os.Getenv("REDDIT_CLIENT_SECRET")
	redditUser         string  = os.Getenv("REDDIT_USER")
	redditPassword     string  = os.Getenv("REDDIT_PASSWORD")
	lcolor             *bool   = hermes.CLIFlags.Bool("C", true, "enable colored output")
	llvl               *int    = hermes.CLIFlags.Int("V", 4, "verbosity level (1..6)")
	lfmt               *string = hermes.CLIFlags.String("log", "cli", "logging/output format ("+strings.Join(logfmts(), "|")+")")
	log                *msg.Logger
)

var (
	timeFilter *time.Duration = hermes.CLIFlags.Duration("L", time.Duration(3)*time.Hour, "X(ns, us, ms, s, m, h) to look back in time for posts")
	post       *bool          = hermes.CLIFlags.Bool("post", false, "send posts to rocketchat channel/user")
)

func init() {
	hermes.CLIFlags.Parse(os.Args[1:])
	var err error
	if log, err = msg.New(
		msg.Formats[*lfmt].String(),
		msg.Formats[msg.CLITimeFmt].String(),
		"subreddit-feed",
		*lcolor,
		msg.Lvl(*llvl),
	); err != nil {
		panic(err)
	}
	if *post {
		if hermes.RCChannel == "" && hermes.RCRoomID == "" {
			log.Panic("need a RoomID or Channel to post the message in..")
		}
		if hermes.RCToken == "" || hermes.RCTokenID == "" {
			log.Panic("need a token/ID to run..")
		}
		if hermes.RCURL == "" {
			log.Panic("need a target RC instance..")
		}
	}
}

func logfmts() []string {
	var logfmts []string
	for f := range msg.Formats {
		if !strings.Contains(f, "rfc") { // exclude the time formats
			logfmts = append(logfmts, f)
		}
	}
	return logfmts
}

func main() {
	s, err := newSession()
	if err != nil {
		log.Panic(err.Error())
	}
	for _, subreddit := range hermes.CLIFlags.Args() {
		log.Info("processing /r/" + subreddit)
		getSubmissions(subreddit, *timeFilter, s)
	}
}

func newSession() (*geddit.OAuthSession, error) {
	s, err := geddit.NewOAuthSession(
		redditClientID,
		redditClientSecret,
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0",
		"",
	)
	if err != nil {
		return nil, err
	}
	err = s.LoginAuth(redditUser, redditPassword)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getSubmissions(subreddit string, add time.Duration, s *geddit.OAuthSession) error {
	r, err := s.AboutSubreddit(subreddit)
	if err != nil {
		return err
	}
	posts, err := s.SubredditSubmissions(
		subreddit,
		geddit.NewSubmissions,
		geddit.ListingOptions{},
	)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	good := 0
	tpl := hermes.POSTTemplate{
		Avatar: r.HeaderImg,
		Text:   "Last " + timeFilter.String() + " posts in " + " */r/" + r.Name + "* " + r.String(),
	}
	attachedPosts := make([]hermes.TemplateAttachments, 0)
	for i := range posts {
		x := posts[i]
		posted := float64Time(x.DateCreated).UTC()
		if posted.Add(add).After(now) {
			log.Noticef(
				"\n---\n%s\n%s\n%s\n%s\n---\n",
				posted.String(),
				x.String(),
				x.URL,
				x.ThumbnailURL,
			)
			if *post {
				attachedPosts = append(
					attachedPosts,
					hermes.TemplateAttachments{
						AuthorName:  x.Author,
						AuthorLink:  "https://www.reddit.com/user/" + x.Author + "/",
						Title:       x.Title,
						TitleLink:   x.FullPermalink(),
						Text:        x.String(),
						MessageLink: x.URL,
						ThumbURL:    x.ThumbnailURL,
						Timestamp:   posted,
					})
			}
			good++
		}
	}
	log.Infof("processed %d posts from /r/%s, %d were in the desired timeframe\n", len(posts), subreddit, good)
	if *post && good > 0 {
		tpl.Attachments = attachedPosts[:]
		log.Noticef("sending posts to RocketChat %s", hermes.RCChannel)
		if status, err := hermes.POSTMessage(tpl); err != nil {
			log.Errorf("%s, status code is %d", err.Error(), status)
		}
	}
	return nil
}

func float64Time(ts float64) time.Time {
	return time.Unix(
		int64(ts),
		int64((ts-float64(int64(ts)))*1e9),
	)
}

func timeFloat64(ts time.Time) float64 {
	return float64(ts.UnixNano()) / 1e9
}
