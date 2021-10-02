build:
	docker build -t financial-nginx-dev .

run: build
	docker run -d \
		--rm \
		--net host \
		financial-nginx-dev
