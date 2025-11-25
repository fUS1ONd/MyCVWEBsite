FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /usr/local/src/app .
COPY --from=builder /usr/local/src/config ./config
COPY --from=builder /usr/local/src/migrations ./migrations

CMD ["./app"]
