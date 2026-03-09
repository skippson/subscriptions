FROM golang:1.26.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o subscriptions ./cmd/app/main.go

FROM alpine:3.19 as runner

WORKDIR /app

RUN adduser -D subscriptions

COPY --from=builder /app/subscriptions .

USER subscriptions

EXPOSE 8888

CMD ["./subscriptions"]