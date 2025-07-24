package main

import (
	"log"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	
	"github.com/alaw22/blog_aggregator/internal/config"
	"github.com/alaw22/blog_aggregator/internal/database"
)


func main(){
	configObj, err := config.Read()

	if err != nil {
		log.Fatal(err)
	}

	// connection to postgres
	db, err := sql.Open("postgres",configObj.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	State := state{
		db:dbQueries,
		conf:&configObj,
	}

	Commands := commands{
		cmds:make(map[string]func(*state, command) error),
	}

	Commands.register("login",handlerLogin)
	Commands.register("register",handlerRegister)
	Commands.register("reset",handlerReset)
	Commands.register("users",handlerShowUsers)
	Commands.register("agg",handlerAgg)
	Commands.register("addfeed",middlewareLoggedIn(handlerAddFeed))
	Commands.register("feeds",handlerShowFeeds)
	Commands.register("follow",middlewareLoggedIn(handlerFollow))
	Commands.register("following",middlewareLoggedIn(handlerShowFollowing))
	Commands.register("unfollow",middlewareLoggedIn(handlerUnfollow))
	


	args := os.Args

	if len(args) < 2 {
		log.Fatal("No command")
	}

	Command := command{
		name: args[1],
		args: args[2:],
	}

	err = Commands.run(&State,Command)
	if err != nil {
		log.Fatal(err)
	}

}
