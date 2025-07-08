# ---- Build Stage ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy
COPY . .
RUN go build -o /auditlog cmd/serverd/main.go
# Install mockery
RUN go install github.com/vektra/mockery/v2@v2.42.0
# ---- Run Stage ----
FROM golang:1.23-alpine
WORKDIR /app
COPY --from=builder /auditlog /auditlog
COPY --from=builder /go/bin/mockery /usr/local/bin/mockery
RUN apk add --no-cache git
COPY . .
EXPOSE 8080
CMD ["/auditlog"]