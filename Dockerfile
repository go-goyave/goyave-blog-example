FROM golang:1.17-alpine

LABEL maintainer="Jérémy LAMBERT (SystemGlitch) <jeremy.la@outlook.fr>"

RUN apk update && apk upgrade && apk add --no-cache git openssh gcc libc-dev
RUN go get github.com/cespare/reflex

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR /app

EXPOSE 8080

CMD dockerize -wait tcp://mariadb:3306 reflex -s -- sh -c 'go run main.go'