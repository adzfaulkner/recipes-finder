module github.com/adzfaulkner/recipes-finder/finder

go 1.13

require (
	github.com/adzfaulkner/recipes-finder/db v0.0.0
	github.com/adzfaulkner/recipes-finder/es v0.0.0
)

replace github.com/adzfaulkner/recipes-finder/db => ../db

replace github.com/adzfaulkner/recipes-finder/es => ../es
