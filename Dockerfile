FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache --update git

RUN go mod download

RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["./main"]