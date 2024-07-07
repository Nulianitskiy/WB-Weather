FROM golang:1.22-alpine

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN go build -o Weather ./cmd/app

EXPOSE 8080

CMD ["./Weather"]