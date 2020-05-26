package finder

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/adzfaulkner/recipes-finder/db"
)

func TestFindRecipes(t *testing.T) {
	var cases = []struct {
		testname       string
		in             []byte
		expectedTitles []string
	}{
		{
			"findRecipes that feature 400 of tomatoes",
			[]byte(`{"id":"test","query":[{"ingredient":"tomatoes","quantity":400}]}`),
			[]string{"Veggie chilli con carne"},
		},
		{
			"findRecipes that feature 400 of tomato",
			[]byte(`{"id":"test","query":[{"ingredient":"tomato","quantity":400}]}`),
			[]string{
				"Cheesey Chorizo Chicken and Spinach",
				"indian spiced lamb",
				"chicken with smoked paprika and almonds",
				"Goan Fish Curry",
				"tomatoes, eggs and chorizo",
				"Veggie chilli con carne",
			},
		},
		{
			"findRecipes that feature 100 of spinach and mushroom",
			[]byte(`{"id":"test","query":[{"ingredient":"mushroom","quantity":100},{"ingredient":"spinach","quantity":100}]}`),
			[]string{
				"Chicken with creamy wild mushroom and tarragon sauce",
				"Creamy Steak and Spinach",
				"super spedy beef stroganoff",
				"indian spiced lamb",
				"chicken with smoked paprika and almonds",
				"Veggie chilli con carne",
				"Cheesey Chorizo Chicken and Spinach",
				"Coddled eggs with spinach and bacon",
			},
		},
		{
			"findRecipes that feature 400 of chicken and 100 of spinach / mushroom",
			[]byte(`{"id":"test","query":[{"ingredient":"chicken","quantity":400},{"ingredient":"mushroom","quantity":100},{"ingredient":"spinach","quantity":100}]}`),
			[]string{
				"Chicken with creamy wild mushroom and tarragon sauce",
				"chicken with smoked paprika and almonds",
				"Cheesey Chorizo Chicken and Spinach",
				"Creamy Steak and Spinach",
				"super spedy beef stroganoff",
				"indian spiced lamb",
				"Veggie chilli con carne",
				"Coddled eggs with spinach and bacon",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.testname, func(t *testing.T) {
			var d struct {
				ID string `json:"id"`
			}

			json.Unmarshal(c.in, &d)

			repo := db.GetHitsRepo(context.TODO())

			var m Message
			json.Unmarshal(c.in, &m)

			FindRecipes(c.in, &repo, context.TODO())

			doc, _ := repo.Get(context.TODO(), d.ID)

			if len(doc.Hits.Hits) != len(c.expectedTitles) {
				t.Errorf("Actual titles length %d and expected length %d mismatch", len(doc.Hits.Hits), len(c.expectedTitles))
			}

			for i, hit := range doc.Hits.Hits {
				if hit.Source.Title != c.expectedTitles[i] {
					t.Errorf("Actual %sand expected %s mismatch", hit.Source.Title, c.expectedTitles[i])
				}
			}
		})
	}
}
