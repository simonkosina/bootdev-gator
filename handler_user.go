package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("'login' expects a single username argument")
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("'login' failed to set user: %w", err)
	}

	fmt.Printf("User has been set to: %s", cmd.args[0])

	return nil
}
