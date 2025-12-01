# Build time image
FROM golang:1.25.1 AS builder

# Private repos
ENV GOPRIVATE=github.com/nahuelsoma/*
ARG GITHUB_TOKEN

# Set the token to use for private repos
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Create the working directory
WORKDIR /app

# Copy the project files
COPY . .

# Download private dependencies using the token
RUN go mod download

# Build the app
RUN go build -o main .

# Final image
FROM debian:bookworm-slim

# Install SSL certificates for HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/main .

# Run the app
CMD ["./main"]
