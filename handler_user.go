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
		return fmt.Errorf("'login' expects a single username argument\n")
	}

	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("'login' failed to find user: %w\n", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("'login' failed to set user: %w\n", err)
	}

	fmt.Printf("Current user has been set to: %s\n", name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("'register' expects a single username argument\n")
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
		return fmt.Errorf("'register' failed to create user: %w\n", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("'register' failed to set user: %w\n", err)
	}

	fmt.Printf("User was created successfully: %+v\n", user)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("'users' doesn't expect any arguments\n")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("'users' failed to retrieve users: %w\n", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Println(user.Name, "(current)")
		} else {
			fmt.Println(user.Name)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("'reset' doesn't expect any arguments\n")
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("'reset' failed to delete users: %w\n", err)
	}

	fmt.Printf("Users table was successfully reset\n")

	return nil
}
