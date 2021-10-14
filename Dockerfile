# Build UI
FROM node:14-alpine as ui

WORKDIR /ui
COPY ./ui .

RUN yarn
RUN yarn build

# Build API
FROM golang:1.17.2-alpine as api

WORKDIR /api
COPY ./api .
RUN go mod vendor
RUN go mod tidy
RUN go build -o api ./cmd/main.go

# Combine into one image
FROM alpine:3

RUN apk --no-cache add tzdata
WORKDIR /opt
COPY --from=api /api/api /opt
COPY --from=ui /ui/build /opt/build
COPY ./entrypoint.sh /opt

ENTRYPOINT ["/opt/entrypoint.sh"]
