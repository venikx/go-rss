package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/venikx/go-rss/database"

	"context"
	"log"
	"sync"

	"encoding/xml"
	"io"
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

				log.Println("feed", feed.Name)
				_, err = db.MarkFeedFetched(context.Background(), feed.Id)

				if err != nil {
					log.Println("error marking fetch feed: ", err)
					return
				}

				rssFeed, err := urlToFeed(feed.Url)
				if err != nil {
					log.Println("error fetching feeds: ", err)
					return
				}

				for _, item := range rssFeed.Channel.Item {
					log.Println("found post", item.Title, "on feed", feed.Name)
				}

				log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
			}()
		}
		wg.Wait()
	}
}

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	go startScraping(NewServer.db, 10, time.Minute)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
