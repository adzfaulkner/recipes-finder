package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/adzfaulkner/recipes-finder/db"
)

type IngredientResponse struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Measure  string `json:"measure"`
}

type RecipeResponse struct {
	Title       string               `json:"title"`
	Serves      int                  `json:"serves"`
	Ingredients []IngredientResponse `json:"ingredients"`
}

type ResultResponse struct {
	ID             string           `json:"id"`
	Success        bool             `json:"success"`
	GoodResults    []RecipeResponse `json:"good_results"`
	PartialResults []RecipeResponse `json:"partial_results"`
}

func Results(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
		resp := ErrorResponse{
			Success: false,
			Message: "ID is missing from request",
		}

		respJson, _ := json.Marshal(resp)

		handleResponse(w, respJson)
		return
	}

	bCtx := context.Background()
	hitsRepo := db.GetHitsRepo(bCtx)
	defer hitsRepo.Cleanup(bCtx)
	ctx, cancel := context.WithTimeout(bCtx, 15*time.Second)
	defer cancel()
	res, err := hitsRepo.Get(ctx, id[0])

	if err != nil || !res.Found {
		resp := ErrorResponse{
			Success: false,
			Message: "In progress",
		}

		respJson, _ := json.Marshal(resp)

		handleResponse(w, respJson)
		return
	}

	good, partial := generateResults(&res)

	resp := ResultResponse{
		ID:             id[0],
		Success:        true,
		GoodResults:    good,
		PartialResults: partial,
	}

	respJson, _ := json.Marshal(resp)
	handleResponse(w, respJson)
}

func generateResults(r *db.RecipeDocument) ([]RecipeResponse, []RecipeResponse) {
	partial := []RecipeResponse{}
	good := []RecipeResponse{}
	var ingredientResponses []IngredientResponse
	var recipeResponse RecipeResponse

	for _, h := range *&r.Hits.Hits {
		ingredientResponses = []IngredientResponse{}

		for _, i := range h.Source.Ingredients {
			ingredientResponses = append(ingredientResponses, IngredientResponse{
				Item:     i.Item,
				Quantity: i.Quantity,
				Measure:  i.Measure,
			})
		}

		recipeResponse = RecipeResponse{
			Title:       h.Source.Title,
			Serves:      h.Source.Serves,
			Ingredients: ingredientResponses,
		}

		if r.Hits.MaxScore == h.Score {
			good = append(good, recipeResponse)
		} else {
			partial = append(partial, recipeResponse)
		}
	}

	return good, partial
}
