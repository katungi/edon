FROM golang:1.25-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main ./cmd/web

# Copy static files
COPY cmd/web/static ./cmd/web/static

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
