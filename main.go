package main

import (
	"fmt"
	"log"
	"github.com/alaw22/blog_aggregator/internal/config"
)

func main(){
	configObj, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	configObj.SetUser("alex")
	newConfigObj, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("db_url = %s\n",newConfigObj.DB_URL)
	fmt.Printf("current_user_name = %s\n",newConfigObj.CurrentUser)

}
