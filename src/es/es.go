package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

type QueryParam struct {
	Ingredient string
	Quantity   int
}

type QueryParams interface {
	GetQueryParams() []QueryParam
}

func getElasticClient() elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:9200", os.Getenv("ES_HOST")),
		},
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Unable to connect to elastic: %s", err)
	}

	_, err = es.Info()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	return *es
}

func generateQuery(qp *[]QueryParam) bytes.Buffer {
	qb := generateQueryBody(qp)

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": qb,
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(q)

	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	return buf
}

func generateQueryBody(qp *[]QueryParam) *[]map[string]interface{} {
	qrySegs := []map[string]interface{}{}
	for _, p := range *qp {
		qrySegs = append(qrySegs, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "ingredients",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{
								"wildcard": map[string]interface{}{
									"ingredients.item": map[string]interface{}{
										"value": fmt.Sprintf("*%s*", p.Ingredient),
										"boost": 3.0,
									},
								},
							},
						},
						"should": []interface{}{
							map[string]interface{}{
								"range": map[string]interface{}{
									"ingredients.quantity": map[string]interface{}{
										"lte":   p.Quantity,
										"boost": 2.0,
									},
								},
							},
						},
					},
				},
			},
		})
	}

	return &qrySegs
}

func SearchRecipes(qp QueryParams) map[string]interface{} {
	qps := qp.GetQueryParams()

	q := generateQuery(&qps)

	es := getElasticClient()

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("recipes"),
		es.Search.WithBody(&q),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	var r map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&r)

	return r
}
