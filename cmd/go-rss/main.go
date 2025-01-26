package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/venikx/go-rss/database"
	"github.com/venikx/go-rss/server"
)

type RSSFeed struct {
	Channel struct {
		XMLName     xml.Name  `xml:"channel"`
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Item        []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

func urlToFeed(url string) (RSSFeed, error) {
	rssFeed := RSSFeed{}
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return rssFeed, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return rssFeed, err
	}
	fmt.Println(string(data))

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return rssFeed, err
	}

	return rssFeed, nil

}

func startScraping(
	db database.Service,
	concurrency int,
	timeBetweenRequest time.Duration) {
	log.Printf("scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), concurrency)
		if err != nil {
			log.Printf("error fetching feeds %v", err)
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go func() {
				defer wg.Done()

				_, err = db.MarkFeedFetched(context.Background(), feed.Id)

				if err != nil {
					log.Println("error fetching feeds: ", err)
					return
				}

				rssFeed, err := urlToFeed(feed.Url)
				if err != nil {
					log.Println("error fetching feeds: ", err)
					return
				}

				for _, item := range rssFeed.Channel.Item {
					log.Println("found post: ", item.Title)
				}
			}()
		}
		wg.Wait()
	}
}

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}
