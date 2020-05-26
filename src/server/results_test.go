package main

import (
	"context"
	"testing"

	"github.com/adzfaulkner/recipes-finder/db"
	"github.com/adzfaulkner/recipes-finder/finder"
)

var cases = []struct {
	testname               string
	message                []byte
	expectedGoodResults    []string
	expectedPartialResults []string
}{
	{
		testname: "Search and find recipes using a valid request body",
		message:  []byte(`{"id":"test","query":[{"ingredient":"mushroom","quantity":100},{"ingredient":"spinach","quantity":100},{"ingredient":"chicken","quantity":400}]}`),
		expectedGoodResults: []string{
			"Chicken with creamy wild mushroom and tarragon sauce",
		},
		expectedPartialResults: []string{
			"chicken with smoked paprika and almonds",
			"Cheesey Chorizo Chicken and Spinach",
			"Creamy Steak and Spinach",
			"super spedy beef stroganoff",
			"indian spiced lamb",
			"Veggie chilli con carne",
			"Coddled eggs with spinach and bacon",
		},
	},
	{
		testname: "Search and find recipes using a valid request body",
		message:  []byte(`{"id":"test","query":[{"ingredient":"tomato","quantity":500},{"ingredient":"mushroom","quantity":100}]}`),
		expectedGoodResults: []string{
			"Veggie chilli con carne",
		},
		expectedPartialResults: []string{
			"Cheesey Chorizo Chicken and Spinach",
			"super spedy beef stroganoff",
			"indian spiced lamb",
			"chicken with smoked paprika and almonds",
			"Goan Fish Curry",
			"tomatoes, eggs and chorizo",
			"Chicken with creamy wild mushroom and tarragon sauce",
			"Creamy Steak and Spinach",
		},
	},
}

func TestResultsResponse(t *testing.T) {
	for _, c := range cases {
		c := c
		t.Run(c.testname, func(t *testing.T) {
			repo := db.GetHitsRepo(context.TODO())
			finder.FindRecipes(c.message, &repo, context.TODO())

			r := DoResultsRequest(t, "test")

			if r.ID != "test" {
				t.Errorf("Actual response ID %s expected response ID %s result mismatch", r.ID, "test")
			}

			if len(r.GoodResults) != len(c.expectedGoodResults) {
				t.Errorf("Actual good matches length %d and expected good matches length %d mismatch", len(r.GoodResults), len(c.expectedGoodResults))
			}

			if len(r.PartialResults) != len(c.expectedPartialResults) {
				t.Errorf("Actual good matches length %d and expected good matches length %d mismatch", len(r.GoodResults), len(c.expectedPartialResults))
			}

			for i, result := range r.GoodResults {
				if result.Title != c.expectedGoodResults[i] {
					t.Errorf("Actual %s and expected %s mismatch", result.Title, c.expectedGoodResults[i])
				}
			}

			for i, result := range r.PartialResults {
				if result.Title != c.expectedPartialResults[i] {
					t.Errorf("Actual %s and expected %s mismatch", result.Title, c.expectedPartialResults[i])
				}
			}
		})
	}
}
