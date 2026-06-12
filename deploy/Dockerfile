FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/api .
COPY --from=builder /app/database ./database

EXPOSE 8080

CMD ["./api"]