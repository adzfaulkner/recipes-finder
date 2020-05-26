package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecipeDocumentHit struct {
	Score  float64 `bson:"_score"`
	Source struct {
		Title       string `bson:"title"`
		Serves      int    `bson:"serves"`
		Ingredients []struct {
			Item     string `bson:"item"`
			Quantity int    `bson:"quantity"`
			Measure  string `bson:"measure"`
		} `bson:"ingredients"`
	} `bson:"_source"`
}

type RecipeDocumentHits struct {
	MaxScore float64             `bson:"max_score"`
	Hits     []RecipeDocumentHit `bson:"hits"`
}

type RecipeDocument struct {
	Found bool               `bson:"found"`
	Hits  RecipeDocumentHits `bson:"hits"`
}

type RepositoryInterface interface {
	Create(ctx context.Context, id string, found bool, hits interface{})
	Update(ctx context.Context, id string, found bool, hits interface{})
	Get(ctx context.Context, id string) (RecipeDocument, error)
	Delete(ctx context.Context, id string)
	Cleanup(ctx context.Context)
}

type HitsRepository struct {
	Client *mongo.Client
}

func getRecipesCollection(c *mongo.Client) *mongo.Collection {
	return c.Database(os.Getenv("MONGO_DATABASE")).Collection("hits")
}

func (h *HitsRepository) Create(ctx context.Context, id string, found bool, hits interface{}) {
	collection := getRecipesCollection(h.Client)

	_, err := collection.InsertOne(ctx, bson.M{"id": id, "hits": hits, "found": found})

	if err != nil {
		log.Fatalf("Mongo create error: %s", err)
	}
}

func (h *HitsRepository) Update(ctx context.Context, id string, found bool, hits interface{}) {
	collection := getRecipesCollection(h.Client)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.D{
			{"$set", bson.M{"hits": hits, "found": found}},
		},
	)

	if err != nil {
		log.Fatalf("Mongo update error: %s", err)
	}
}

func (h *HitsRepository) Get(ctx context.Context, id string) (RecipeDocument, error) {
	collection := getRecipesCollection(h.Client)

	var res RecipeDocument

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&res)

	return res, err
}

func (h *HitsRepository) Delete(ctx context.Context, id string) {
	collection := getRecipesCollection(h.Client)
	collection.DeleteMany(ctx, bson.M{"id": id})
}

func (h *HitsRepository) Cleanup(ctx context.Context) {
	h.Client.Disconnect(ctx)
}

func GetHitsRepo(ctx context.Context) HitsRepository {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017", os.Getenv("MONGO_ROOT_USERNAME"), os.Getenv("MONGO_ROOT_PASSWORD"), os.Getenv("MONGO_HOST"))))

	if err != nil {
		log.Fatalf("Mongo connect error: %s", err)
	}

	return HitsRepository{
		Client: client,
	}
}
