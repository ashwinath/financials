commit=$(shell git rev-parse HEAD)

build:
	docker build -t ashwinath/financials:$(commit) .
