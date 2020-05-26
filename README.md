# Recipe Finder App

What else you can do whilst in lockdown other than try to extend your skills? Not much. For a long time, I wanted to delve into new tech such as GoLang, VueJS hence I created this app which finds you recipes for the ingredients you input.

## Getting started

This app serves as a protoype therefore should only be installed on one's local machine.

### Prerequisites

1. Docker installed
2. docker-compose installed (optional but recommended for obvious reasons)

### Installing

1. Clone this repo and cd into it.
2. Copy the root .env.sample & paste as ./.env
3. Run `docker network create intersect-backend`. This creates a network so that various container's can communicate to one another.
4. Modify the newly created .env file values as necessary
5. Run `docker-compose up`
6. Once docker has finished bringing up the app, head to http://localhost:8080 in your browser and search away.

### Developing

In the root directory, there is a docker-compose-dev.yml that will provide you with a:

1. A VueJS client environment that you can use to develop with. Reachable on localhost:5000 by default.
2. A mongo-express client which allows you to administer the DB with relative ease. Reachable on localhost:8082 by default.

### Tests

To run the go functional/unit tests, run the following command:

`docker-compose -f ./docker-compose-test.yml up --abort-on-container-exit`

There are some Cypress frontend behavioural tests.

It's not as simple as running the commands as you need to install the dependencies as described here:

https://www.cypress.io/blog/2019/05/02/run-cypress-with-a-single-docker-command/

Once your environment is ready, then simply run either:

`./e2e/cy_open.sh` - This will open up the task runner app in XQuartz

`./e2e/cy_run.sh` - This will execute all tests in terminal.

## TODO

1. Evolve elastic search query in terms of scoring relevance. For example, if you enter "tomato" results for "cherry tomatoes" are returned.
2. Increase test coverage. Start looking at how to mock dependencies instead of relying on the expensive functional tests.
3. Design a better UI. Currently it's simply designed as a test harness.