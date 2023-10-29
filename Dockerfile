FROM golang:alpine

# Download and install dockerize
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz

# Set destination for COPY
WORKDIR /app

RUN apk update && apk add --no-cache git

# Enable Go modules
ENV GO111MODULE=auto

# Download Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/main"]