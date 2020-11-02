# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.12.0-alpine3.9

# Add Maintainer Info
LABEL maintainer="Shreyas Srinivasa <shsr@es.aau.dk>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 2222,8080,23,1883

# Command to run the executable
RUN apt-get update
RUN docker pull eclipse-mosquitto
RUN  docker run -it -p 1883:1883 eclipse-mosquitto
CMD ["./bin/riotpot"]
