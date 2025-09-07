package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/simonkosina/bootdev-gator/internal/config"
	"github.com/simonkosina/bootdev-gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
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

	if len(os.Args) < 2 {
		// TODO: Print proper usage with all the commands
		log.Fatal("Usage: gator <command> [args...]\n")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&st, cmd); err != nil {
		log.Fatal(err)
	}
}
