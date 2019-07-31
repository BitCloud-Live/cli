# This file is a template, and might need editing before it works on your project.
FROM golang:1.10 AS builder
ARG version
WORKDIR $GOPATH/src/github.com/yottab/cli
COPY . ./
WORKDIR $GOPATH/src
RUN CGO_ENABLED=0 go build -a -installsuffix nocgo --tags netgo -ldflags '-w -X github.com/yottab/cli/cmd.version=$(version) -extldflags "-static"' -o /yb -i github.com/yottab/cli/main.go
#RUN go build -v  -o cron1 youtab/cron/expWorker/main.go
#RUN go build -v  -o cron2 youtab/cron/mailWorker/main.go

FROM alpine
RUN apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
ENV HOME /app
ENV PATH ${PATH}:${HOME}
WORKDIR ${HOME}
COPY --from=builder /yb yb
CMD [ "/app/yb" ]