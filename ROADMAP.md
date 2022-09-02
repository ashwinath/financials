# Roadmap

Ideation of what features to add in.

## Querying house price to add to asset value

Consider making multiple calls to get a 1-2 year window of `month` and multiple range of `lease_commence_date`. To research if range queries are allowed.

```bash
curl --silent \
    --GET \
    --data-urlencode 'q={
        "street_name": "changi rd",
        "flat_type": "4 ROOM",
        "month": "2022-01",
        "lease_commence_date": "1977"
    }' \
    --data-urlencode 'resource_id=f1765b54-a209-4718-8d38-a39237f502b3' \
    --data-urlencode 'sort=month desc' \
    'https://data.gov.sg/api/action/datastore_search' | jq '.result.records | join(",")'
```

## Mortgage

- Interest left
- Interest paid
- Principal left
- Principal paid
- Total paid
- Total left
- % of mortgage done
