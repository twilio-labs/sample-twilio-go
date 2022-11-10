FROM golang:1.19.0

# Set the Current Working Directory inside the container
WORKDIR /app

RUN export GO111MODULE=on
RUN export CGO_ENABLED=0
RUN export GOOS=linux
RUN export GOARCH=arm64

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

# App environment variables specific to running from this Docker image
ENV TEMPLATE_DIR=/app/pkg/template ASSET_DIR=/app/asset

# Build the application
RUN mkdir -p out/bin
RUN go build -o out/bin/sample-twilio-go.out cmd/app/main.go

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./out/bin/sample-twilio-go.out"]
