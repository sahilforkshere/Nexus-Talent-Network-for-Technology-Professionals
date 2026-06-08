FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY jobs-svc/go.mod jobs-svc/go.sum ./
ENV GOTOOLCHAIN=auto
RUN GOWORK=off go mod download
COPY jobs-svc/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o jobs-svc .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/jobs-svc .
EXPOSE 4002
CMD ["./jobs-svc"]
