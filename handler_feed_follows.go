package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonkosina/bootdev-gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: gator %s <feed_url>\n", cmd.name)
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: gator %s\n", cmd.name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("'following' failed to retrieve followed feeds: %w\n", err)
	}

	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: gator %s <feed_url>\n", cmd.name)
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("'unfollow' failed to find feed: %w\n", err)
	}

	err = s.db.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("'unfollow' failed to unfollow: %w\n", err)
	}

	fmt.Printf("'%s' unfollowed '%s' feed\n", user.Name, feed.Name)

	return nil
}
