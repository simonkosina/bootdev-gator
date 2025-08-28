package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonkosina/bootdev-blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("'login' expects a single username argument")
	}

	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("'login' failed to find user: %w", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("'login' failed to set user: %w", err)
	}

	fmt.Printf("Current user has been set to: %s\n", name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("'register' expects a single username argument")
	}

	id := uuid.New()
	currentTime := time.Now().UTC()
	name := cmd.args[0]

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        id,
			UpdatedAt: currentTime,
			CreatedAt: currentTime,
			Name:      name,
		})

	if err != nil {
		return fmt.Errorf("'register' failed to create user: %v", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("'register' failed to set user: %w", err)
	}

	fmt.Printf("User was created successfully: %v", user)

	return nil
}
