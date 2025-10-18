FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migration ./cmd/migration


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migration .

EXPOSE 8080

CMD ["/app/main"]