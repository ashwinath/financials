test:
	@go test -v -race -short -cover -coverprofile cover.out ./...
	@go tool cover -func cover.out

setup:
	@docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password --name financial postgres:13-alpine
