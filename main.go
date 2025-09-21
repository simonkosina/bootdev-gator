package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/simonkosina/bootdev-gator/internal/config"
	"github.com/simonkosina/bootdev-gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func printHelp() {
	helpText := `Usage: gator <command> [args...]

Commands:
  login <user_name>                Set the current user
  register <user_name>             Register a new user and set as current
  reset                            Reset (delete) all users
  users                            List all users
  agg <time_between_reqs>          Collect feeds every given duration (e.g. 1m, 30s)
  addfeed <feed_name> <feed_url>   Add a new feed and follow it
  feeds                            List all feeds
  follow <feed_url>                Follow a feed
  following                        List feeds you are following
  unfollow <feed_url>              Unfollow a feed
  browse [limit]                   Browse posts (default limit: 2)
`
	fmt.Print(helpText)
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	st := state{
		cfg: &cfg,
		db:  database.New(db),
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&st, cmd); err != nil {
		if strings.Contains(err.Error(), "Unknown command:") {
			log.Print(err)
			printHelp()
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}
}
