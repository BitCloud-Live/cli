# This file is a template, and might need editing before it works on your project.
FROM golang:1.10 AS builder
WORKDIR $GOPATH/src/github.com/uvcloud/uv-cli
# This will download deps in docker file ignored for faster build(copy vendor folder)
# ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
# RUN chmod +x /usr/bin/dep
# COPY Gopkg.toml Gopkg.lock ./
# RUN dep ensure --vendor-only
COPY . ./
WORKDIR $GOPATH/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o uv-cli github.com/uvcloud/uv-cli/main.go

FROM alpine:latest
WORKDIR /usr/bin/
COPY --from=builder /go/src/uv-cli uv-cli

