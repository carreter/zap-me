FROM golang:1.22

WORKDIR /app

# Get dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY *.go ./
RUN go build -o /zap-me-server

# Run the server
EXPOSE 9000
CMD ["/zap-me-server"]