package main

import (
	"fmt"
	"context"
	"strconv"

	"github.com/alaw22/blog_aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var(
		err error 
		limit int
	)

	limit = 2

	if len(cmd.args) != 0{
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Need post limit to be of type INT: %w\n",err)
		}
	}

	postsForUserParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(),
									   postsForUserParams)

    if err != nil {
		return err
	}

	fmt.Printf("Here are the %d most recent posts:\n\n",limit)

	
	for _, post := range posts {
		fmt.Printf("Title: %s\nPublished: %v\nURL: %s\nDescription:\n\t%s\n",
			post.Title,
			post.PublishedAt,
			post.Url,
			post.Description.String,
		)
		
		fmt.Printf("\n----------------------------------------\n")

	}

	return nil
	
}