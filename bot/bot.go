package bot

import (
	"net/http"
	"net/url"
)

const botApiUrl string = "https://api.telegram.org/bot"

type (
	Bot struct {
		Token string
	}
)

func (bot *Bot) execApi(method string, values string) {
	_, err := http.Get(botApiUrl + bot.Token + "/" + method + "?" + values)
	if err != nil {
		return
	}
}

func (bot *Bot) SendMessage(chat_id string, text string, options map[string]string) {
	values := url.Values{}
	values.Set("chat_id", chat_id)
	values.Set("text", text)
	for option, optionVal := range options {
		values.Set(option, optionVal)
	}
	bot.execApi("sendMessage", values.Encode())
}
