package main

import (
	"fmt"
	"time"
	"errors"
	"context"
	"github.com/google/uuid"

	"github.com/alaw22/blog_aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected an URL to follow")
	}

	url := cmd.args[0]

	// look up if the feed already exists if not then we can't perform the operation
	feed, err := s.db.GetFeedByURL(context.Background(),url)
	if err != nil {
		return fmt.Errorf("Feed doesn't exist to follow: %w\n",err)
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID, 
	}

	_, err = s.db.CreateFeedFollow(context.Background(),createFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println("You are now following",feed.Name)

	return nil
}

func handlerShowFollowing(s *state, cmd command, user database.User) error {

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("You are following these feeds:\n")
	for _, feedFollow := range feedFollows {
		fmt.Printf(" - %s\n",feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected a url")
	}

	url := cmd.args[0]

	unfollowParams := database.UnfollowParams{
		UserID: user.ID,
		Url: url,
	}

	err := s.db.Unfollow(context.Background(), unfollowParams)
	if err != nil {
		return fmt.Errorf("Unable to unfollow feed: %w\n",err)
	}

	fmt.Println("Successfully unfollowed feed:",url)

	return nil

}