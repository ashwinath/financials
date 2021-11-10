#!/bin/bash
set -e

jsonnet -J vendor financials.jsonnet > dashboard-temp.json

payload="{\"dashboard\": $(jq . dashboard-temp.json), \"overwrite\": true}"

curl -X POST $BASIC_AUTH \
  -H 'Content-Type: application/json' \
  -d "${payload}" \
  "http://admin:admin@localhost:3000/api/dashboards/db"

rm dashboard-temp.json
