commit=$(shell git rev-parse HEAD)
all: build push

build:
	docker build -t $(REGISTRY)/financials:$(commit) .

push:
	docker push $(REGISTRY)/financials:$(commit)

setup:
	@docker run -d \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=password \
		--name financial \
		postgres:13-alpine
	@timeout 30 bash -c "until docker exec financial pg_isready; do sleep 2; done"

migrate:
	@docker run \
		--net="container:financial" \
		--volume ${PWD}/migrations:/etc/migrations \
		migrate/migrate:v4.15.0 \
		-path=/etc/migrations \
		-database='postgres://postgres:password@localhost:5432/postgres?sslmode=disable' up

deps:
	go mod vendor
	go mod tidy

test:
	DB_HOST="localhost" \
	DB_PORT="5432" \
	DB_USER="postgres" \
	DB_NAME="postgres" \
	DB_PASSWORD="password" \
	DB_TIMEZONE="Asia/Singapore" \
		go test -race -v -short -cover -coverprofile cover.out ./...
	go tool cover -func cover.out
