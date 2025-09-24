package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/edgarmueller/go-gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	user, _ := s.db.GetUser(context.Background(), name)
	if user.Name == name {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("couldn't register user: %w", err)
	}

	s.cfg.SetUser(user.Name)

	fmt.Printf("User %s registered successfully!\n", user.Name)
	fmt.Printf("User ID: %s\n", user.ID)
	return nil
}
