FROM golang:1.22 as builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o certificate-converter .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /srv/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/certificate-converter .

# Run the binary
CMD ["./certificate-converter"]
