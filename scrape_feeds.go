package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/edgarmueller/go-gator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	now := time.Now()
	feed, err := s.db.GetNextFeedToFetch(context.Background(), sql.NullTime{
		Time:  now,
		Valid: true,
	})

	if err != nil {
		return err
	}

	s.db.MarkFeedAsFetched(context.Background(), database.MarkFeedAsFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt:     now,
	})

	rssFeed, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v\n", feed.Name)

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Title: %+v\n", item.Title)
		savePosts(s, feed.ID, rssFeed)
	}

	return nil
}

func savePosts(s *state, feedID uuid.UUID, rssFeed *RSSFeed) error {
	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			publishedAt = time.Now()
		}

		newPost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			FeedID:      uuid.NullUUID{UUID: feedID, Valid: true},
			PublishedAt: publishedAt,
		}

		_, err = s.db.CreatePost(context.Background(), newPost)

		if err != nil {

			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				continue
			}

			return fmt.Errorf("couldn't create post: %w", err)
		}
	}
	return nil
}
