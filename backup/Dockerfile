# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./

# Build - no CGO needed (no SQLite)
RUN CGO_ENABLED=0 GOOS=linux go build -o backup .

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/backup .

# Run as non-root user
RUN adduser -D -u 1000 appuser
USER appuser

CMD ["./backup"]
