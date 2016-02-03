package main

import (
	"encoding/json"
	"io/ioutil"
)

type (
	configuration struct {
		Token    string     `json:"token"`
		Channels []*Channel `json:"channels"`
	}
	Channel struct {
		Name  string  `json:"name"`
		Feeds []*Feed `json:"feeds"`
	}
	Feed struct {
		Url                   string `json:"url"`
		Update                int    `json:"update,string,omitempty"`
		DisableWebPagePreview string `json:"disable_web_page_preview,omitempty"`
		ParseMode             string `json:"parse_mode,omitempty"`
		Template              string `json:"template,omitempty"`
	}
)

func (config *configuration) LoadFromFile() (err error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		return
	}
	return
}
