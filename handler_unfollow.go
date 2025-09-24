package main

import (
	"context"
	"fmt"

	"github.com/edgarmueller/go-gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)

	if err != nil {
		return fmt.Errorf("couldn't get feed by URL: %w", err)
	}

	s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	fmt.Printf("Feed %s is now unfollowed by user %s\n", feed.Name, user.Name)

	return nil
}
