# Start from the official Go image to build our application.
FROM golang:1.21.5 as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the go.mod and go.sum to download all dependencies.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container.
COPY cmd/ cmd/
COPY terminal/ terminal/

# Set the working directory to the cmd directory where the main.go file is located.
WORKDIR /app/cmd

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goaiterminal-chat

# Start a new stage from scratch for a smaller image.
FROM alpine:latest  

# Install ca-certificates in case the application makes HTTPS requests.
RUN apk --no-cache add ca-certificates

# Create a non-root user and switch to it.
RUN adduser -D gogenaiterminal
USER gogenaiterminal

# Set the working directory to the user's home directory.
WORKDIR /home/gogenaiterminal

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/cmd/. /usr/local/bin/

# Run the binary.
ENTRYPOINT ["gogenaiterminal-chat"]
