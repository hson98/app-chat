# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Sơn Trần <tranhuyson41198@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY .. .
#COPY ./ /app
#COPY /docs .
# Build the Go app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
RUN mkdir -p ~/uploads
RUN mkdir -p ~/docs
# Copy the Pre-built binary file from the previous stage. Observe we also copied the app.env file
COPY --from=builder /app/main .
#COPY --from=builder /app/.env .
COPY --from=builder /app/app.env .
COPY --from=builder /app/docs/ ./docs
COPY --from=builder /app/config/ ./config

# Expose port 8080 to the outside world
EXPOSE 8080
#Command to run the executable
CMD ["./main"]