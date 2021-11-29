# financials

![Version: 1.2.1](https://img.shields.io/badge/Version-1.2.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Helm chart for financials

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami | postgresql | 10.12.4 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| financials.assets | string | `"date,type,amount\n2020-03-31,CPF,1000\n2020-03-31,Bank,20000\n2020-03-31,Mortgage,-40000\n2020-03-31,Investments,20000"` | CSV values for the assets |
| financials.assetsCSVDirectory | string | `"/etc/assets"` | folder containing the assets directories |
| financials.config | object | `{"AlphaVantageAPIKey":"test","AssetsCSVFile":"/etc/assets/assets.csv","Database":{"BatchInsertSize":100,"Host":"<release name>-postgresql.<namespace>.svc.cluster.local","Name":"postgres","Password":"password","Port":5432,"TimeZone":"Asia/Singapore","User":"postgres"},"ExpensesCSVFile":"/etc/expenses/expenses.csv","IncomeCSVFile":"/etc/income/income.csv","TradesCSVFile":"/etc/trades/trades.csv"}` | Config for financials api, see api/config.yaml |
| financials.cronSchedule | string | `"0 */4 * * *"` | cron schedule |
| financials.expenses | string | `"date,type,amount\n2020-03-31,Credit Card,500\n2020-03-31,Reimbursement,-200\n2020-03-31,Tithe,800"` | CSV values for the expenses |
| financials.expensesCSVDirectory | string | `"/etc/expenses"` | folder containing the expenses directories |
| financials.financialsGrafana.apiKey | string | `"secret"` |  |
| financials.financialsGrafana.endpoint | string | `"www.example.com:3000"` | URL and port of the grafana site |
| financials.financialsGrafana.image.pullPolicy | string | `"Always"` | Image pull policy in Kubernetes |
| financials.financialsGrafana.image.repository | string | `"ghcr.io/ashwinath/financials-grafana"` | Respository of the image. |
| financials.financialsGrafana.image.tag | string | `"latest"` | Override this value for the desired image tag |
| financials.image.pullPolicy | string | `"Always"` | Image pull policy in Kubernetes |
| financials.image.repository | string | `"ghcr.io/ashwinath/financials"` | Respository of the image. |
| financials.image.tag | string | `"latest"` | Override this value for the desired image tag |
| financials.income | string | `"date,type,amount\n2021-03-11,Base,500\n2021-03-11,Bonus,200"` | CSV values for the income |
| financials.incomeCSVDirectory | string | `"/etc/income"` | folder containing the income directories |
| financials.migrate | object | `{"image":{"tag":"v4.15.0"}}` |  choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, uncomment the following lines, adjust them as necessary, and remove the curly braces after 'resources:'. limits:   cpu: 100m   memory: 128Mi requests:   cpu: 100m   memory: 128Mi |
| financials.migrate.image.tag | string | `"v4.15.0"` | migrate/migrate image tag for db migrations |
| financials.resources | object | `{}` | Resources requests and limits for the financial app |
| financials.tradeCSVDirectory | string | `"/etc/trades"` | folder containing the csv files |
| financials.trades | string | `"date_purchased,symbol,trade_type,price_each,quantity\n2021-03-11,IWDA.LON,buy,76.34,10"` | CSV values for the trades |
| postgresql.persistence | object | `{"enabled":true}` |  choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, uncomment the following lines, adjust them as necessary, and remove the curly braces after 'resources:'. limits:   cpu: 100m   memory: 128Mi requests:   cpu: 100m   memory: 128Mi |
| postgresql.persistence.enabled | bool | `true` | Persist Postgresql data in a Persistent Volume Claim  |
| postgresql.postgresqlDatabase | string | `"postgres"` | Database name for Turing Postgresql database |
| postgresql.postgresqlPassword | string | `"password"` | Password for postgresql database, highly recommended to change this value |
| postgresql.postgresqlUsername | string | `"postgres"` | Username for postgresql database |
| postgresql.resources | object | `{}` | Resources requests and limits for the database |

