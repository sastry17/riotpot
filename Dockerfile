FROM golang:1.13.5 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  ./bin/riotpot ./cmd/riotpot

FROM alpine:latest AS prod

COPY --from=builder /app .

EXPOSE 2222

EXPOSE 23

CMD ["./bin/riotpot"]