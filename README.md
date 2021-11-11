# Financials

Financials is a way to track your financial independence track in Singapore's context. This is still very much a work in progress, as of now, only the investments portions work. We use Grafana as the frontend.

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
