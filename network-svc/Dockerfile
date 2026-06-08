FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY network-svc/go.mod network-svc/go.sum ./
ENV GOTOOLCHAIN=auto
RUN GOWORK=off go mod download
COPY network-svc/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o network-svc .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/network-svc .
EXPOSE 4003
CMD ["./network-svc"]
