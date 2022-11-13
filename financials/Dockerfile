FROM rust:1.65.0-slim as builder
LABEL org.opencontainers.image.authors="ashwinath@hotmail.com"
LABEL org.opencontainers.image.source https://github.com/ashwinath/financials

WORKDIR /usr/src/financials
COPY . .
RUN apt update && apt install -y libssl-dev libpq-dev pkg-config
RUN cargo install --path .

FROM debian:bullseye-slim
WORKDIR /app
RUN apt update && apt install -y libssl-dev libpq-dev ca-certificates
COPY --from=builder /usr/local/cargo/bin/financials /app/financials
COPY ./migrations ./migrations
CMD ["./financials"]
