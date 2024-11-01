FROM golang:1.22.2-alpine3.19

RUN mkdir -p /home/app

COPY . /home/app

WORKDIR /home/app

RUN go mod tidy

RUN go build ./cmd/server/main.go

CMD [ "./main"]