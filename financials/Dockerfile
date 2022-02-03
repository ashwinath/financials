FROM rust:1.58.1-slim as builder
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
