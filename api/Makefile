setup:
	@docker stop financial
	@docker rm financial
	@docker run -d \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=password \
		--name financial \
		postgres:13-alpine
	@sleep 10
	@docker run \
		--net="container:financial" \
		--volume ${PWD}/migrations:/etc/migrations \
		migrate/migrate:v4.15.0 \
		-path=/etc/migrations \
		-database='postgres://postgres:password@localhost:5432/postgres?sslmode=disable' up

deps:
	go mod vendor
	go mod tidy

test: setup deps
	DB_HOST="localhost" \
	DB_PORT="5432" \
	DB_USER="postgres" \
	DB_NAME="postgres" \
	DB_PASSWORD="password" \
	DB_TIMEZONE="Asia/Singapore" \
		go test -v -race -short -cover -coverprofile cover.out ./...
	go tool cover -func cover.out
