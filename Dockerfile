FROM golang:1.19.0

# Set the Current Working Directory inside the container
WORKDIR /app

RUN export GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

# Build the application
RUN make build-app

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./out/bin/webinar-scale-up-app.out"]