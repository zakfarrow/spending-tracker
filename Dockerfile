FROM golang:1.23-alpine

# Install build dependencies
RUN apk add --no-cache git

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Run templ generate to generate initial templates
RUN templ generate

# Expose port 8080
EXPOSE 8080

# Use air for hot reloading
CMD ["air"]
