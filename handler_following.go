package main

import (
	"context"
	"fmt"

	"github.com/edgarmueller/go-gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsByUserID(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Printf("User %s is not following any feeds.\n", user.Name)
		return nil
	}

	fmt.Printf("Feeds followed by user %s:\n", user.Name)
	for _, ff := range feedFollows {
		if err != nil {
			return fmt.Errorf("couldn't get feed by ID: %w", err)
		}
		fmt.Printf("- %s (%s)\n", ff.FeedName, ff.FeedUrl)
	}

	fmt.Printf("Total feeds followed: %d\n", len(feedFollows))
	return nil
}
