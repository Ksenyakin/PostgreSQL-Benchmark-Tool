FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]
