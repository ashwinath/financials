# Financials

Financials is a way to track your financial independence track in Singapore's context. This is still very much a work in progress, as of now, only the investments portions work.

## Installing

Helm is the preferred way to install this. First, get an API key here: https://www.alphavantage.co/support/#api-key

Then deploy using this command.

```bash
helm repo add ashwinath https://ashwinath.github.io/helm-charts/
helm repo update
helm upgrade financials ashwinath/financials \
    --install \
    --wait \
    --namespace=financials \
    --set financials.config.Database.Host="financials-postgresql.financials.svc.cluster.local" \
    --set financials.config.AlphaVantageAPIKey="<Your alphavantage key here>" \
    --set financials.image.repository="ghcr.io/ashwinath/financials" \
    --set financials.image.tag="bdbeff842b44fa038b49ab089650817ab8043e8d" 
