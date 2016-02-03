package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shpaker/rsschanbot/bot"
	"github.com/shpaker/rsschanbot/rss"
)

const (
	DEFAULT_UPDATE_TIME int    = 5
	DEFAULT_TEMPLATE    string = "%_ITEM_TITLE_%\n%_ITEM_LINK_%"
	CONFIGURATION_FILE  string = "./config.json"
)

var (
	err          error
	outLastItems int
	telegramBot  *bot.Bot
	config       *configuration
)

func init() {
	flag.IntVar(&outLastItems, "o", 0, "Print {count} of last items in RSS channel")
	defer flag.Parse()

	config = &configuration{}
	if err = config.LoadFromFile(); err != nil {
		log.Panic(err)
	}
	if config.Token == "" {
		log.Fatal("Empty token")
	}
	telegramBot = &bot.Bot{config.Token}
}

func startFeedListen(chatId string, feed *Feed) {
	rssFeed := rss.NewRss(feed.Url)
	if feed.Update <= 0 {
		feed.Update = DEFAULT_UPDATE_TIME
	}
	if feed.Template == "" {
		feed.Template = DEFAULT_TEMPLATE
	}
	rssChan, err := rssFeed.NewUpdateRssChan(time.Minute * time.Duration(feed.Update))
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case outLastItems < -1:
		outLastItems = 0
	case outLastItems == -1 || outLastItems > len(rssFeed.Channel.Items):
		outLastItems = len(rssFeed.Channel.Items)
	}

	for i := 0; i < outLastItems; i++ {
		r := strings.NewReplacer(
			"%_ITEM_TITLE_%", rssFeed.Channel.Items[i].Title,
			"%_ITEM_LINK_%", rssFeed.Channel.Items[i].Link,
			"%_ITEM_DESCRIPTION_%", rssFeed.Channel.Items[i].Description)
		post := r.Replace(feed.Template)
		telegramBot.SendMessage(chatId, post, nil)
	}

	for {
		select {
		case msg := <-rssChan:
			r := strings.NewReplacer(
				"%_ITEM_TITLE_%", msg.Title,
				"%_ITEM_LINK_%", msg.Link,
				"%_ITEM_DESCRIPTION_%", msg.Description)
			post := r.Replace(feed.Template)
			telegramBot.SendMessage(chatId, post, nil)
		}
	}
}

func main() {
	for _, channel := range config.Channels {
		for _, feed := range channel.Feeds {
			go startFeedListen(channel.Name, feed)
		}
	}
	var input string
	fmt.Scanln(&input)
}
