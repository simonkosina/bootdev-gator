package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonkosina/bootdev-blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("'addFeed' expects name and url arguments\n")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	currentTime := time.Now().UTC()

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("'addFeed' failed to retrieve current user: %w\n", err)
	}

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
		return fmt.Errorf("'addFeed' failed to create feed: %w\n", err)
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
		return fmt.Errorf("'addFeed' failed to create feed follow: %w\n", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("'feeds' doesn't expect any arguments\n")
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
