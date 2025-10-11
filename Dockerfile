# Stage 1: Build the application
FROM golang:1.25.2-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's build cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# -ldflags="-s -w" reduces binary size by omitting debugging information and symbol table
# -o app specifies the output binary name
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

# Stage 2: Create a minimal runtime image
FROM alpine:3.22.2

WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/app .

# Expose the port your Gin application listens on (e.g., 8080)
EXPOSE 8080

# Command to run the application when the container starts
CMD ["./app"]