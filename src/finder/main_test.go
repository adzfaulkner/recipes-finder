package finder

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/adzfaulkner/recipes-finder/db"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	hitsRepo := db.GetHitsRepo(context.TODO())
	hitsRepo.Create(context.TODO(), "test", false, nil)

	fmt.Printf("\033[1;36m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	hitsRepo := db.GetHitsRepo(context.TODO())
	hitsRepo.Delete(context.TODO(), "test")
	hitsRepo.Cleanup(context.TODO())

	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}
