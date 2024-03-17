FROM golang:1.22

WORKDIR /usr/src/app

COPY ./GouelServer/ ./
COPY ./release/env/server.env .env

RUN go mod download && go mod verify

CMD [ "go", "run", "main.go"]