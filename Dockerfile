# -------------------------------------
# Stage 1: Build Stage
# -------------------------------------
FROM golang:1.25.3-alpine3.22 AS builder

# Set working directory
WORKDIR /app

# Install git for go modules and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o finance-tracker main.go

# -------------------------------------
# Stage 2: Final Stage
# -------------------------------------
FROM alpine:3.22

# Install CA certificates for HTTPS requests
# RUN apk add --no-cache ca-certificates

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser

# Set timezone (optional, but fixes 'unknown time zone' issue)
ENV TZ=Asia/Kolkata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/finance-tracker .

# Copy environment file
COPY .env .env

# Copy migrations file
COPY migrations migrations

# Ensure correct permissions
RUN chown -R appuser:appuser /app

# Run as non-root for security
USER appuser

# # Expose the port the app will run on
# EXPOSE 8000

# Set default command
CMD ["./finance-tracker"]
