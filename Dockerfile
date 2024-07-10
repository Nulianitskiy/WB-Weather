# Builder stage
FROM golang:1.22-alpine as builder

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN go build -o Weather ./cmd/app

# Distribution stage
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

EXPOSE 8080

COPY --from=builder /go/src/app/Weather /app/
COPY .env /app/.env

CMD ["./Weather"]