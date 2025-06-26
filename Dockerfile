# Step 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go app (static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o mailer

# Step 2: Use a minimal image to run it
FROM alpine:3.20

WORKDIR /app

# Copy only the built binary
COPY --from=builder /app/mailer .

# Expose the app port
EXPOSE 8888

# Run the binary
ENTRYPOINT ["./mailer"]
