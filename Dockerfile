FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheme-manager-adapter .
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/scheme-manager-adapter .
COPY --from=builder /app/config.yaml ./
EXPOSE 8080
CMD ["./scheme-manager-adapter"]