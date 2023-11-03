ARG RUST_VERSION
ARG DEBIAN_VERSION
FROM rust:${RUST_VERSION}-slim as builder
LABEL org.opencontainers.image.authors="ashwinath@hotmail.com" \
      org.opencontainers.image.source="https://github.com/ashwinath/financials"

WORKDIR /usr/src/financials
COPY . .
RUN apt update && apt install -y libssl-dev libpq-dev pkg-config
RUN cargo install --path .

FROM debian:${DEBIAN_VERSION}
WORKDIR /app
RUN apt update && apt install -y libssl-dev libpq-dev ca-certificates
COPY --from=builder /usr/local/cargo/bin/financials /app/financials
CMD ["./financials"]
