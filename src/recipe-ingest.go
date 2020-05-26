package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Ingredient struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Measure  string `json:"measure"`
}

type Recipe struct {
	Title       string       `json:"title"`
	Serves      int          `json:"serves"`
	Ingredients []Ingredient `json:"ingredients"`
}

func main() {
	csvFile, _ := os.Open("./recipes.csv")
	f, _ := os.Create("./recipes.json")
	f.Truncate(0)

	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	var serves int
	var recipe Recipe

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if line[1] == "" {
			continue
		}

		serves, _ = strconv.Atoi(line[0])
		recipe = Recipe{
			Title:       line[1],
			Serves:      serves,
			Ingredients: findIngredients(&line),
		}

		recipeJson, _ := json.Marshal(recipe)

		f.WriteString("{ \"index\":{} }\n")
		f.Write(recipeJson)
		f.WriteString("\n")
	}
}

func findIngredients(line *[]string) []Ingredient {
	var ingredients []Ingredient
	re := regexp.MustCompile("^([\\d]+)([\\w]{0,})")
	var res []string
	var quantity int
	var measure string

	for i := 2; i < len(*line); i = i + 2 {
		l := *line

		res = re.FindStringSubmatch(l[i+1])

		if len(res) < 2 {
			continue
		}

		quantity, _ = strconv.Atoi(res[1])
		measure = res[2]

		ingredients = append(ingredients, Ingredient{
			Item:     l[i],
			Quantity: quantity,
			Measure:  measure,
		})
	}

	return ingredients
}
