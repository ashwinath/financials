commit=$(shell git rev-parse HEAD)
all: build push

build:
	docker build -t $(REGISTRY)/financials-grafana:$(commit) -t $(REGISTRY)/financials-grafana:latest .

push:
	docker push $(REGISTRY)/financials-grafana:$(commit)
	docker push $(REGISTRY)/financials-grafana:latest

grafana:
	docker run -d --net host -e GF_INSTALL_PLUGINS=grafana-piechart-panel --name grafana grafana/grafana

load:
	bash load.sh
