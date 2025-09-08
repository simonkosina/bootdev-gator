package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/simonkosina/bootdev-gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: gator %s <feed_name> <feed_url>\n", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]
	currentTime := time.Now().UTC()

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("'addfeed' failed to create feed: %w\n", err)
	}

	fmt.Printf("Feed was added successfully: %+v\n", feed)

	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("'addfeed' failed to create feed follow: %w\n", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: gator %s\n", cmd.name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("'feeds' failed to retrieve feeds: %w\n", err)
	}

	for _, feed := range feeds {
		fmt.Printf("%+v\n", feed)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: gator %s <time_between_reqs>\n", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("'agg' error parsing time duration (expected Go's time.ParseDuration format, e.g. 1s, 1m, 1h30m): %w\n", err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("'agg' time between requests must be greater than 0\n")
	}

	ticker := time.NewTicker(timeBetweenRequests)

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
