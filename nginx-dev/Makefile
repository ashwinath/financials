build:
	docker build -t financial-nginx-dev .

run: build
	docker run --rm \
		--net host \
		financial-nginx-dev
