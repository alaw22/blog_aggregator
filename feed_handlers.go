package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"
	
	"github.com/alaw22/blog_aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	const testURL = "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), testURL)
	if err != nil {
		return err
	}

	fmt.Printf("Title:%s\nLink:%s\nDescription:%s\n",
			   rssFeed.Channel.Title,
			   rssFeed.Channel.Link,
			   rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {

		fmt.Printf("[%02d]\nTitle:%s\nLink:%s\nDescription:%s\n",
				i,
				item.Title,
				item.Link,
				item.Description)
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