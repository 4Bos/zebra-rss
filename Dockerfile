FROM golang:1.25.3-alpine3.22

WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY backend/go.mod ./
RUN go mod download

# Copy source code
COPY backend .

# Expose port
EXPOSE 8080

# Use air for hot reload
CMD ["air", "-c", ".air.toml"]