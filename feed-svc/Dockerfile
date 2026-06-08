FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY feed-svc/go.mod feed-svc/go.sum ./
ENV GOTOOLCHAIN=auto
RUN GOWORK=off go mod download
COPY feed-svc/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o feed-svc .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/feed-svc .
EXPOSE 4004
CMD ["./feed-svc"]
