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

	id := uuid.New()
	name := cmd.args[0]
	url := cmd.args[1]
	currentTime := time.Now().UTC()

	userName := s.cfg.CurrentUserName
	if len(userName) == 0 {
		return fmt.Errorf("'addFeed' requires a user to be logged in\n")
	}

	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("'addFeed' failed to retrieve current user: %w\n", err)
	}

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        id,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)

	fmt.Printf("Feed was added successfully: %+v\n", feed)

	return nil
}
