# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache make gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/bin/api .
#COPY --from=builder /app/migrations ./migrations
#COPY --from=builder /app/prisma ./prisma

# Install Node.js and npm for Prisma
#RUN apk add --no-cache nodejs npm

# Install Prisma CLI
#RUN npm install -g prisma

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./api"]
