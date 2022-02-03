commit=$(shell git rev-parse HEAD)
all: build push

build:
	docker build -t $(REGISTRY)/financials:$(commit) -t $(REGISTRY)/financials:latest .

push:
	docker push $(REGISTRY)/financials:$(commit)
	docker push $(REGISTRY)/financials:latest

db:
	@docker run -d \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=password \
		--name financial \
		postgres:13-alpine
	@timeout 30 bash -c "until docker exec financial pg_isready; do sleep 2; done"
