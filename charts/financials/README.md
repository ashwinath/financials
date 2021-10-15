# financials

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.1](https://img.shields.io/badge/AppVersion-0.0.1-informational?style=flat-square)

Helm chart for financials

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami | postgresql | 10.12.4 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| financials.config | object | `{"AlphaVantageAPIKey":"test","Database":{"BatchInsertSize":100,"Host":"<release name>-postgresql.<namespace>.svc.cluster.local","Name":"postgres","Password":"password","Port":5432,"TimeZone":"Asia/Singapore","User":"postgres"},"PortfolioCalculationInterval":"1h","Server":{"Port":8000,"ReactFilePath":"../ui/build","ReadTimeoutInSeconds":"5s","WriteTimeoutInSeconds":"5s"}}` | Config for financials api, see api/config.yaml |
| financials.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy in Kubernetes |
| financials.image.repository | string | `"ghcr.io/ashwinath/financials"` | Respository of the image. |
| financials.image.tag | string | `"latest"` | Override this value for the desired image tag |
| financials.migrate | object | `{"image":{"tag":"v4.15.0"}}` |  choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, uncomment the following lines, adjust them as necessary, and remove the curly braces after 'resources:'. limits:   cpu: 100m   memory: 128Mi requests:   cpu: 100m   memory: 128Mi |
| financials.migrate.image.tag | string | `"v4.15.0"` | migrate/migrate image tag for db migrations |
| financials.podAnnotations | object | `{}` | Kubernetes pod annotations in key pair value |
| financials.replicaCount | int | `1` | Number of replicas |
| financials.resources | object | `{}` | Resources requests and limits for the financial app |
| financials.service.port | int | `80` | Kubernetes service port |
| financials.service.type | string | `"ClusterIP"` | Kubernetes service type |
| postgresql.persistence | object | `{"enabled":true}` |  choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, uncomment the following lines, adjust them as necessary, and remove the curly braces after 'resources:'. limits:   cpu: 100m   memory: 128Mi requests:   cpu: 100m   memory: 128Mi |
| postgresql.persistence.enabled | bool | `true` | Persist Postgresql data in a Persistent Volume Claim  |
| postgresql.postgresqlDatabase | string | `"postgres"` | Database name for Turing Postgresql database |
| postgresql.postgresqlPassword | string | `"password"` | Password for postgresql database, highly recommended to change this value |
| postgresql.postgresqlUsername | string | `"postgres"` | Username for postgresql database |
| postgresql.resources | object | `{}` | Resources requests and limits for the database |
