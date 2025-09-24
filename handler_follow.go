package main

import (
	"context"
	"fmt"
	"time"

	"github.com/edgarmueller/go-gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)

	if err != nil {
		return fmt.Errorf("couldn't get feed by URL: %w", err)
	}

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	fmt.Printf("Feed %s is now followed by user %s\n", feed.Name, user.Name)

	return nil
}
