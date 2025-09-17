package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/simonkosina/bootdev-gator/internal/database"
)

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

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("Usage: gator %s <limit>\n", cmd.name)
	}

	var limit int32

	if len(cmd.args) == 0 {
		log.Println("No limit argument provided for 'browse', defaulting to 2")
		limit = 2
	} else {
		res, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("'browse' couldnt parse provided limit ('%s') to integer: %w\n", cmd.args[0], err)
		}
		limit = int32(res)
	}

	posts, err := s.db.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("'browse' couldn't retrieve posts for user '%s': %w\n", user.Name, err)
	}

	const (
		bold      = "\033[1m"
		underline = "\033[4m"
		blue      = "\033[34m"
		green     = "\033[32m"
		cyan      = "\033[36m"
		reset     = "\033[0m"
	)

	log.Printf("Found %d posts for user '%s'\n", len(posts), user.Name)
	for i, post := range posts {
		fmt.Printf("%s: %s%s%s\n", post.FeedName, bold, post.Title, reset)
		fmt.Printf("  %sURL:%s %s%s%s\n", cyan, reset, blue, post.Url, reset)
		fmt.Printf("  %sPublished At:%s %s\n", cyan, reset,
			post.PublishedAt.Local().Format("Mon, 02 Jan 2006 15:04:05"))

		if post.Description.Valid && post.Description.String != "" {
			fmt.Printf("  %sDescription:%s %s\n", cyan, reset, post.Description.String)
		}

		if i < len(posts)-1 {
			fmt.Println(" ", green+"─────────────────────────────────────────────"+reset)
		}
	}

	return nil
}
