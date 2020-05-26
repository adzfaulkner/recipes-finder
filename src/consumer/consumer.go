package main

import (
	"context"
	"log"

	"github.com/adzfaulkner/recipes-finder/db"
	"github.com/adzfaulkner/recipes-finder/finder"
	"github.com/adzfaulkner/recipes-finder/rabbit"
)

func main() {
	client := rabbit.GetClient("recipes")
	defer client.Cleanup()

	msgs := client.Consume()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			ctx := context.Background()
			hitsRepo := db.GetHitsRepo(ctx)
			finder.FindRecipes(d.Body, &hitsRepo, ctx)
			hitsRepo.Cleanup(ctx)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
