package main

import (
	"fmt"
	"time"
	"context"
	"errors"
	"database/sql"
	"github.com/google/uuid"
	
	"github.com/alaw22/blog_aggregator/internal/database"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error in GetNextFeedToFetch(): %w\n",err)
	}

	markFeedFetchedParams := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
		ID: feed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(),markFeedFetchedParams)	
	if err != nil {
		return fmt.Errorf("Couldn't mark feed '%s' as fetched: %w\n",feed.Name,err)
	}

	fmt.Printf("Feed '%s' was successfully fetched and marked as such\n\n",feed.Name)

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Unable to fetch feed '%s': %w\n",feed.Name,err)
	}
	
	fmt.Printf("Title: %s\nLink: %s\nDescription: %s\n\n",
			   rssFeed.Channel.Title,
			   rssFeed.Channel.Link,
			   rssFeed.Channel.Description)
	
	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("[%03d] Title: %s\n", i, item.Title)
	}


	fmt.Print("\n-----------------------------------------------\n")



	return nil


}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected time_between_requests argument")
	}

	time_string := cmd.args[0]

	time_between_requests, err := time.ParseDuration(time_string)
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_requests)

	ticker := time.NewTicker(time_between_requests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough arguments for addfeed\n")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: name,
		Url: url,
		UserID: user.ID,
	}

	_, err := s.db.CreateFeed(context.Background(),feedParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed entry made for",url)

	// Follow said feed
	follow_cmd := command{
		name: "follow",
		args: []string{url},
	}

	err = handlerFollow(s, follow_cmd, user)
	if err != nil {
		fmt.Printf("Unable to follow added feed: %w\n",err)
	}

	return nil
}

func handlerShowFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		username, err := s.db.GetUserNameFromID(context.Background(),feed.UserID)
		if err != nil {
			return fmt.Errorf("No user with same id as entry: %w\n",err)
		}
		fmt.Printf("Name: %s\nURL: %s\nUsername: %s\n",
					feed.Name,
					feed.Url,
					username)
	}

	return nil
}