FROM docker.io/bitnami/minideb:buster

COPY --from=bitnami/jsonnet:0.17.0 /opt/bitnami/jsonnet/bin/jsonnet /usr/bin/jsonnet
RUN apt update && apt install -y curl bash jq
RUN curl -L https://github.com/jsonnet-bundler/jsonnet-bundler/releases/download/v0.4.0/jb-linux-amd64 -o /usr/bin/jb
RUN chmod +x /usr/bin/jb

WORKDIR app
COPY . .
RUN jb install
