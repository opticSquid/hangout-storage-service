FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Copy only necessary files
COPY go.mod go.sum ./
RUN go mod download
COPY . . 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o storage-service .

# Create a new stage for the final image
FROM ubuntu:noble

# Install necessary packages with optimizations
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates ffmpeg gpac && apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy the application.yaml from resources directory from the build context to final image
COPY --from=builder /app/resources/application.yaml /resources/application.yaml
# Copy the built binary from the previous stage
COPY --from=builder /app/storage-service /

# Set the working directory
WORKDIR /

# Run binary
CMD ["./storage-service"]