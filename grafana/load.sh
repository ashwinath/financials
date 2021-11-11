#!/bin/bash
set -e

jsonnet -J vendor financials.jsonnet > dashboard-temp.json

payload="{\"dashboard\": $(jq . dashboard-temp.json), \"overwrite\": true}"

curl -X POST $BASIC_AUTH \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer ${GRAFANA_API_KEY}" \
  -d "${payload}" \
  "http://${GRAFANA_ENDPOINT:-localhost:3000}/api/dashboards/db"

rm dashboard-temp.json
