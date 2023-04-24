commit=$(shell git rev-parse HEAD)
all: build push

test:
	cargo test

build:
	docker build -t $(REGISTRY)/financials:$(commit) -t $(REGISTRY)/financials:latest .

push:
	docker push $(REGISTRY)/financials:$(commit)
	docker push $(REGISTRY)/financials:latest

db:
	@docker run -d \
		-p 6000:5432 \
		-e POSTGRES_PASSWORD=password \
		--name financial \
		postgres:13-alpine
	@timeout 30 bash -c "until docker exec financial pg_isready; do sleep 2; done"

run:
	@cargo run -- \
		--database-url 'postgres://postgres:password@localhost:6000/postgres' \
		--alphavantage-key ${ALPHAVANTAGE_KEY} \
		--trades-csv sample/trades.csv \
		--assets-csv sample/assets.csv \
		--income-csv sample/income.csv \
		--expenses-csv sample/expenses.csv \
		--mortgage-yaml sample/mortgage.yaml \
		--shared-expense-csv sample/shared_expenses.csv \
		--telegram-bot-endpoint 'http://localhost:8080'
