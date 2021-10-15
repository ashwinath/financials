commit=$(shell git rev-parse HEAD)

all: build push

build:
	docker build -t $(REGISTRY)/financials:$(commit) .

push:
	docker push $(REGISTRY)/financials:$(commit)
