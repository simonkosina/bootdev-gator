package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/simonkosina/bootdev-blog-aggregator/internal/config"
	"github.com/simonkosina/bootdev-blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
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

	if len(os.Args) < 2 {
		// TODO: Print proper usage with all the commands
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&st, cmd); err != nil {
		log.Fatal(err)
	}
}
