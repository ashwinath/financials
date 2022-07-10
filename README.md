# Financials

Financials is a way to track your financial independence track in Singapore's context. This is still very much a work in progress, see the 'Feature Roadmap' for more details. We use Grafana as the frontend and it is assumed that Grafana is installed at the moment.

![](./img/sample-screenshot.png)

## Feature roadmap

This list is not exhaustive but a wishlist that I would work on when I'm free. Mostly credit to this [link](https://www.reddit.com/r/singaporefi/comments/p9p668/the_vital_ratios_to_track_on_your_journey_to/).

- [x] Stock Portfolio
- [x] Monthly expenditure
- [x] Income
- [x] Total Net assets
- [x] Automatic detection of investment amount (right now it is manually keyed in by user)
- [ ] Liabilities
- [ ] Financial Independence Quotient


## Installing

Helm is the preferred way to install this. First, get an API key here: https://www.alphavantage.co/support/#api-key

Then deploy using this command, see `charts/financials/values.yaml` for the defaults.

```bash
helm repo add ashwinath https://ashwinath.github.io/helm-charts/
helm repo update
helm upgrade financials ashwinath/financials \
    --install \
    --wait \
    --namespace=financials \
    --values=values.yaml
```

A sample `values.yaml` may look something like this, with all the trade information inside.

```yaml
financials:
  alphavantageAPIKey: yourAlphaVantageAPIKeyHere

  financialsGrafana:
    apiKey: yourAPIKeyHere
    # or your grafana endpoint here.
    endpoint: grafana.infra

  trades: |-
    date_purchased,symbol,trade_type,price_each,quantity
    2021-03-11,IWDA.LON,buy,76.34,250

  assets: |-
    date,type,amount
    2020-03-31,CPF,1000
    2020-03-31,Bank,20000
    2020-03-31,Investments,20000

  expenses: |-
    date,type,amount
    2020-03-31,Credit Card,500
    2020-03-31,Reimbursement,-200
    2020-03-31,Tithe,800

  income: |-
    date,type,amount
    2021-03-11,Base,500
    2021-03-11,CPF,500
    2021-03-11,Bonus,200
    2021-03-11,CPF Bonus,200

postgresql:
  postgresqlPassword: somePasswordHere
  persistence:
    storageClass: openebs-hostpath
    size: 10Gi
```
