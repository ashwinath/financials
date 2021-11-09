# Build API
FROM golang:1.17.3-alpine as app

WORKDIR /app
COPY . .
RUN go mod vendor
RUN go mod tidy
RUN go build -o app ./cmd/main.go

# Combine into one image
FROM alpine:3

RUN apk --no-cache add tzdata
WORKDIR /opt
COPY --from=app /app/app /opt
COPY ./entrypoint.sh /opt

ENTRYPOINT ["/opt/entrypoint.sh"]
