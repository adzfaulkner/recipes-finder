package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adzfaulkner/recipes-finder/db"
	"github.com/adzfaulkner/recipes-finder/rabbit"
	"github.com/lithammer/shortuuid"
)

type Query struct {
	Ingredient string `json:"ingredient"`
	Quantity   int    `json:"quantity"`
}

type RabbitMessage struct {
	ID    string  `json:"id"`
	Query []Query `json:"query"`
}

type SearchResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == "options" {
		handleResponse(w, []byte{})
		return
	}

	client := rabbit.GetClient("recipes")

	defer client.Cleanup()

	var qry []Query
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&qry)

	if err != nil {
		resp := ErrorResponse{
			Success: false,
			Message: "Unexpected post body received",
		}

		respJson, _ := json.Marshal(resp)

		handleResponse(w, respJson)
		return
	}

	id := shortuuid.New()

	bCtx := context.Background()
	hitsRepo := db.GetHitsRepo(bCtx)
	defer hitsRepo.Cleanup(bCtx)
	ctx, cancel := context.WithTimeout(bCtx, 15*time.Second)
	defer cancel()
	hitsRepo.Create(ctx, id, false, nil)

	body := RabbitMessage{
		ID:    id,
		Query: qry,
	}

	msg, _ := json.Marshal(&body)

	client.Publish(msg)

	resp := SearchResponse{
		ID:      id,
		Success: true,
		Message: "OK",
	}

	respJson, _ := json.Marshal(resp)
	handleResponse(w, respJson)
}

func handleResponse(w http.ResponseWriter, r []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(r)

	if err != nil {
		log.Fatalf("Error writing response header: %s", err)
	}
}
