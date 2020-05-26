module github.com/adzfaulkner/recipes-finder

go 1.13

replace github.com/adzfaulkner/recipes-finder/rabbit => ./rabbit

replace github.com/adzfaulkner/recipes-finder/db => ./db

replace github.com/adzfaulkner/recipes-finder/es => ./es

replace github.com/adzfaulkner/recipes-finder/finder => ./finder

require (
	github.com/adzfaulkner/recipes-finder/db v0.0.0
	github.com/adzfaulkner/recipes-finder/es v0.0.0
	github.com/adzfaulkner/recipes-finder/finder v0.0.0
	github.com/adzfaulkner/recipes-finder/rabbit v0.0.0
	github.com/google/uuid v1.1.1 // indirect
	github.com/lithammer/shortuuid v3.0.0+incompatible
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71 // indirect
)
