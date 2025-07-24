package main

import (
	"fmt"
	"time"
	"errors"
	"context"
	"github.com/google/uuid"
	
	"github.com/alaw22/blog_aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected username")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(),username)
	if err != nil {
		return fmt.Errorf("User '%s' doesn't exist: %w\n",username,err)
	}

	err = s.conf.SetUser(user.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected username")
	}

	username := cmd.args[0]
	
	user, err := s.db.GetUser(context.Background(),username)
	if err == nil {
		return fmt.Errorf("User '%s' already exists\n",username)
	}

	userParams := database.CreateUserParams{
		ID:uuid.New(),
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
		Name: username,
	}


	user, err = s.db.CreateUser(context.Background(),userParams)
	if err != nil {
		return err
	}

	err = s.conf.SetUser(user.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Reset unsuccessful: %w\n",err)
	}

	fmt.Println("Reset successful")

	return nil
}

func handlerShowUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.conf.CurrentUser {
			fmt.Printf("* %s (current)\n",user.Name)
 		} else {
			fmt.Printf("* %s\n",user.Name)
		}
	}

	return nil

}
