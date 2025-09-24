package main

import (
	"context"
	"fmt"
)

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}

	currUser := s.cfg.CurrentUserName

	for _, user := range users {
		name := user.Name
		if name == currUser {
			name += " (current)"
		}
		fmt.Println("* " + name)
	}
	return nil
}
