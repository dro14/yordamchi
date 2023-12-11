# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# The image contains all the necessary build tools and libraries needed to compile Go applications.
FROM golang:1.21 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container to cache the dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main

# Use the Debian based minimal image to run the application.
# This image is extremely small and contains only the bare essentials to run a Go application.
FROM gcr.io/distroless/base-debian10

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /main
COPY --from=builder /app/images.png /images.png

# Expose port 8000 to the outside world
EXPOSE 8000

# Run the web service on container startup.
CMD ["/main"]
