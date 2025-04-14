FROM golang:latest AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# Try to download dependencies, if go.mod exists
RUN if [ -f go.mod ]; then go mod download; fi

# Copy the source code
COPY . .

# Initialize go.mod if it doesn't exist
RUN if [ ! -f go.mod ]; then go mod init universal_api; fi

# Install dependencies
RUN go get github.com/gin-gonic/gin
RUN go get github.com/PuerkitoBio/goquery@v1.8.1
RUN go get gorm.io/gorm
RUN go get gorm.io/driver/sqlite

# Build the application
RUN go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
