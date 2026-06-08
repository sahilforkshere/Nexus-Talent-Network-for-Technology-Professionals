FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY profile-svc/go.mod profile-svc/go.sum ./
ENV GOTOOLCHAIN=auto
RUN GOWORK=off go mod download
COPY profile-svc/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o profile-svc .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/profile-svc .
EXPOSE 4001
CMD ["./profile-svc"]
