#!/bin/sh
echo "Bundling template"

curl -XPUT --header 'Content-Type: application/json' http://localhost:9200/recipes --data-binary @/mappings.json
curl -s -XPOST -H 'Content-Type: application/x-ndjson' http://localhost:9200/recipes/_bulk --data-binary @/recipes.json
