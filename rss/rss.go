package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Rss struct {
		url     string   // ToDo del
		Version string   `xml:"version,attr"`
		Channel *Channel `xml:"channel"`
	}

	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`

		Items []*Item `xml:"item"`

		Language       string `xml:"language"`
		Copyright      string `xml:"copyright"`
		ManagingEditor string `xml:"managingEditor"`
		WebMaster      string `xml:"webMaster"`
		PubDate        string `xml:"pubDate"`
		LastBuildDate  string `xml:"lastBuildDate"`
		Category       string `xml:"category"`
		Generator      string `xml:"generator"`
		Docs           string `xml:"docs"`
		Cloud          string `xml:"cloud"`
		Ttl            string `xml:"ttl"`
		Image          *Image `xml:"image"`
		Rating         string `xml:"rating"`
		TextInput      string `xml:"textInput"`
		SkipHours      string `xml:"skipHours"`
		SkipDays       string `xml:"skipDays"`
	}

	Image struct {
		Url   string `xml:"url"`
		Title string `xml:"title"`
		Link  string `xml:"link"`

		Width       string `xml:"width"`
		Height      string `xml:"height"`
		Description string `xml:"description"`
	}

	Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`

		Author     string   `xml:"author"`
		Categories []string `xml:"category"`
		Comments   string   `xml:"comments"`
		Enclosure  string   `xml:"enclosure"`
		Guid       string   `xml:"guid"`
		PubDate    string   `xml:"pubDate"`
		Source     string   `xml:"source"`
	}
)

func NewRss(url string) (rss *Rss) {
	rss = &Rss{url: url}
	rss.getRss()

	return rss
}

func (rss *Rss) NewUpdateRssChan(updateDuration time.Duration) (itemChan chan *Item, err error) {

	itemChan = make(chan *Item)

	go func() {
		for {
			oldItems := rss.Channel.Items
			if err := rss.getRss(); err != nil {
				continue
			}

			// fmt.Printf("Succesfull updated\n")

			for _, loadedItem := range rss.Channel.Items {
				isNewItem := true
				for _, oldItem := range oldItems {
					if loadedItem.Link == oldItem.Link {
						isNewItem = false
					}
				}
				if isNewItem {
					itemChan <- loadedItem
				}
			}
			time.Sleep(updateDuration)
		}
	}()

	return itemChan, nil
}

func (rss *Rss) getRss() (err error) {

	res, err := http.Get(rss.url)
	if err != nil {
		return
	}

	resData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	err = xml.Unmarshal(resData, rss)
	if err != nil {
		return
	}
	return
}
