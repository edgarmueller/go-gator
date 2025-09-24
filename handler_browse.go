package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/edgarmueller/go-gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {

	var err error
	var limit int

	if len(cmd.Args) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit argument: %v", err)
		}
	}

	posts, err := s.db.GetPostsByUserID(context.Background(), database.GetPostsByUserIDParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	for _, post := range posts {
		// determine desired output width (override with GATOR_WIDTH env var)
		width := 80
		if v := os.Getenv("GATOR_WIDTH"); v != "" {
			if w, err := strconv.Atoi(v); err == nil && w >= 40 {
				width = w
			}
		}

		sep := strings.Repeat("-", width)

		fmt.Println(sep)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Link:  %s\n", post.Url)
		fmt.Printf("Published: %s\n\n", post.PublishedAt)

		desc := strings.TrimSpace(htmlToText(post.Description.String))
		if desc == "" {
			fmt.Println("  (no description)")
			fmt.Println()
		} else {
			fmt.Print(wrapAndIndent(desc, width, "  "))
			fmt.Print("\n")
		}

		fmt.Println(sep)
		fmt.Println()
	}

	return nil
}
