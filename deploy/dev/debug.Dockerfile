FROM golang:1.16.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.16.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN go get github.com/go-delve/delve/cmd/dlv@v1.6.0
RUN mkdir -p /application
ADD . /application
WORKDIR /application

ARG bin_name=app
ENV BIN_NAME=$bin_name

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ./bin/${BIN_NAME} cmd/${BIN_NAME}/main.go

FROM debian:buster-slim
RUN set -xe && apt-get update && apt-get install -y curl

ARG bin_name=app
ENV BIN_NAME=$bin_name

COPY --from=builder /application/bin/${BIN_NAME} /app

EXPOSE 8080 40000
WORKDIR /

COPY --from=builder /go/bin/dlv /
CMD ["/dlv", "--continue", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app"]