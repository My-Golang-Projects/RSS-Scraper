package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/My-Golang-Projects/RSS-Scraper/internal/database"
	"github.com/google/uuid"
)

// this is a long running job so this will run in background
// concurreny how many goroutines do we want the scraping on
func startScraping(db *database.Queries, concurreny int, timeBetweenRequest time.Duration) {
	log.Print("Scraping on %v goroutines every %s duration", concurreny, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	// execute when a new value comes to the ticker channel
	// doing ;; will run immediately and wait for the ticker,
	// using range would have waited a minute
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurreny))
		if err != nil {
			log.Println("error fetcing feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error feching feed:", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID},
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("failed to create post: ", err)
		}
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}

	log.Printf("Feed %s collected, %v oists found", feed.Name, len(rssFeed.Channel.Item))
}
