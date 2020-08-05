# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.14-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Derevtsov Konstantin <godzie@yandex.ru>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/lw/main.go

EXPOSE 8080

CMD ["./main"]
