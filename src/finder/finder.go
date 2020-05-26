package finder

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/adzfaulkner/recipes-finder/db"
	"github.com/adzfaulkner/recipes-finder/es"
)

type MessageQuery struct {
	Ingredient string `json:"ingredient"`
	Quantity   int    `json:"quantity"`
}

type Message struct {
	ID    string         `json:"id"`
	Query []MessageQuery `json:"query"`
}

func (m *Message) GetQueryParams() []es.QueryParam {
	ret := []es.QueryParam{}

	for _, l := range m.Query {
		ret = append(ret, es.QueryParam{
			Ingredient: l.Ingredient,
			Quantity:   l.Quantity,
		})
	}

	return ret
}

func FindRecipes(b []byte, hitsRepo db.RepositoryInterface, ctx context.Context) {
	var m Message
	err := json.Unmarshal(b, &m)

	if err != nil {
		log.Fatalf("Unable to decode message: %s", err)
	}

	match(m, hitsRepo, ctx)
}

func match(m Message, hitsRepo db.RepositoryInterface, ctx context.Context) {
	res := es.SearchRecipes(&m)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	hitsRepo.Update(ctx, m.ID, true, res["hits"])
}
