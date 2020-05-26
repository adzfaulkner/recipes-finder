package es

import (
	"encoding/json"
	"strings"
	"testing"
)

var cases = []struct {
	testname                  string
	in                        []QueryParam
	expectedGenerateQuery     string
	expectedGenerateQueryBody string
}{
	{
		"test 1 query",
		[]QueryParam{{
			Ingredient: "test",
			Quantity:   1,
		}},
		`{"query":{"bool":{"should":[{"nested":{"path":"ingredients","query":{"bool":{"must":[{"wildcard":{"ingredients.item":{"boost":3,"value":"*test*"}}}],"should":[{"range":{"ingredients.quantity":{"boost":2,"lte":1}}}]}}}}]}}}`,
		`[{"nested":{"path":"ingredients","query":{"bool":{"must":[{"wildcard":{"ingredients.item":{"boost":3,"value":"*test*"}}}],"should":[{"range":{"ingredients.quantity":{"boost":2,"lte":1}}}]}}}}]`,
	},
	{
		"test 0 queries",
		[]QueryParam{},
		`{"query":{"bool":{"should":[]}}}`,
		"[]",
	},
}

func TestGenerateQuery(t *testing.T) {
	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.testname, func(t *testing.T) {
			t.Parallel()
			r := generateQuery(&c.in)
			a := strings.TrimRight(r.String(), "\n")

			if strings.Compare(a, c.expectedGenerateQuery) != 0 {
				t.Errorf("Actual %s and expected %s mismatch", a, c.expectedGenerateQuery)
			}
		})
	}
}

func TestGenerateQueryBody(t *testing.T) {
	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.testname, func(t *testing.T) {
			t.Parallel()
			r := generateQueryBody(&c.in)
			a, _ := json.Marshal(r)

			if string(a) != c.expectedGenerateQueryBody {
				t.Errorf("Actual %s and expected %s mismatch", a, c.expectedGenerateQueryBody)
			}
		})
	}
}
