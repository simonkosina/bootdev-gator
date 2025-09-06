package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonkosina/bootdev-blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: gator %s <feed_url>\n", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("'follow' failed to retrieve current user: %w\n", err)
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("'follow' failed to find feed: %w\n", err)
	}

	currentTime := time.Now().UTC()

	follow, err := s.db.CreateFeedFollow(
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
		return fmt.Errorf("'follow' failed to create feed follow: %w\n", err)
	}

	fmt.Printf("'%s' now follow '%s' feed\n", follow.UserName, follow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: gator %s\n", cmd.name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("'following' failed to retrieve followed feeds: %w\n", err)
	}

	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}

	return nil
}
