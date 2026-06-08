FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY search-svc/go.mod search-svc/go.sum ./
ENV GOTOOLCHAIN=auto
RUN GOWORK=off go mod download
COPY search-svc/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o search-svc .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/search-svc .
EXPOSE 4005
CMD ["./search-svc"]
