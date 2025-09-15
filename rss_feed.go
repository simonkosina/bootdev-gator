package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/simonkosina/bootdev-gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	log.Printf("Scraping feed '%s'\n", feed.Name)
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch feed '%s': %v", feed.Name, err)
	}

	_, err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed '%s' as fetched: %v", feed.Name, err)
	}

	for _, item := range feedData.Channel.Item {
		savePost(db, feed, item)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func savePost(db *database.Queries, feed database.Feed, post RSSItem) {
	log.Printf("Saving post '%s' (%s)\n", post.Title, post.Link)

	publishedAt, err := time.Parse(time.RFC1123Z, post.PubDate)
	if err != nil {
		log.Printf("Error parsing 'pubDate' for post '%s': %v\n", post.Title, err)
		return
	}

	currentTime := time.Now().UTC()

	_, err = db.CreatePost(context.Background(), database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Title:     post.Title,
		Url:       post.Link,
		Description: sql.NullString{
			String: post.Description,
			Valid:  len(post.Description) > 0,
		},
		PublishedAt: publishedAt,
		FeedID:      feed.ID,
	})
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
			log.Printf("Post '%s' is already saved, skipping\n", post.Title)
		} else {
			log.Printf("Couldn't save post '%s': %v\n", post.Title, err)
		}

		return
	}
	log.Printf("Post '%s' was saved successfully\n", post.Title)
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feed to fetch", err)
		return
	}

	scrapeFeed(s.db, feed)
}
