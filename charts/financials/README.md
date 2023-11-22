# financials

![Version: 2.3.0](https://img.shields.io/badge/Version-2.3.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Helm chart for financials

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami | postgresql | 13.2.15 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| financials.alphavantageApiKey | string | `"secret"` |  |
| financials.assets | string | `"date,type,amount\n2020-03-31,CPF,1000\n2020-03-31,Bank,20000\n2020-03-31,Mortgage,-40000\n2020-03-31,Investments,20000"` | CSV values for the assets |
| financials.cronSchedule | string | `"0 */4 * * *"` | cron schedule |
| financials.expenses | string | `"date,type,amount\n2020-03-31,Credit Card,500\n2020-03-31,Reimbursement,-200\n2020-03-31,Tithe,800"` | CSV values for the expenses |
| financials.financialsGrafana.apiKey | string | `"secret"` |  |
| financials.financialsGrafana.endpoint | string | `"www.example.com:3000"` | URL and port of the grafana site |
| financials.financialsGrafana.image.pullPolicy | string | `"Always"` | Image pull policy in Kubernetes |
| financials.financialsGrafana.image.repository | string | `"ghcr.io/ashwinath/financials-grafana"` | Respository of the image. |
| financials.financialsGrafana.image.tag | string | `"latest"` | Override this value for the desired image tag |
| financials.image.pullPolicy | string | `"Always"` | Image pull policy in Kubernetes |
| financials.image.repository | string | `"ghcr.io/ashwinath/financials"` | Respository of the image. |
| financials.image.tag | string | `"latest"` | Override this value for the desired image tag |
| financials.income | string | `"date,type,amount\n2021-03-11,Base,500\n2021-03-11,Bonus,200"` | CSV values for the income |
| financials.mortgage | string | `"mortgages:\n- total: 50000.0\n  mortgage_first_payment: 2022-10-10\n  mortgage_duration_in_years: 25\n  mortgage_end_date: 2047-10-10\n  interest_rate_percentage: 2.6\n  downpayments:\n  - date: 2021-10-10\n    sum: 1000.0\n  - date: 2021-12-12\n    sum: 20000.0"` | YAML values for mortgage |
| financials.resources | object | `{}` | Resources requests and limits for the financial app |
| financials.shared_expenses | string | `"date,type,amount\n2023-01-01,Special:Renovations,5000.00\n2023-01-01,Electricity,100.00\n2023-01-01,Water,50.00\n2023-01-01,Gas,30.00\n2023-01-01,Grocery,300.00\n2023-01-01,Eating Out,500.00"` | CSV values for shared expenses |
| financials.telegramBotUrl | string | `"http://<url here>"` | telegram bot dump endpoint, include scheme as well. |
| financials.trades | string | `"date_purchased,symbol,trade_type,price_each,quantity\n2021-03-11,IWDA.LON,buy,76.34,10"` | CSV values for the trades |
| postgresql.auth.postgresqlPassword | string | `"password"` | Password for postgresql database, highly recommended to change this value |
| postgresql.primary.persistence.enabled | bool | `true` | Persist Postgresql data in a Persistent Volume Claim  |
| postgresql.resources | object | `{}` | Resources requests and limits for the database |
