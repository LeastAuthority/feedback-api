# Start from the official Go image based on Alpine
FROM golang:1.19.3-alpine

# Inject the source code
COPY *.go go.* /src

# Download deps, build the app and cleanup the source
WORKDIR /src
RUN go mod download && \
    go mod verify && \
    go build -o /app/feedback-http-server -buildvcs=false -v

# Switch to app directory
WORKDIR /app

# Default smtp config
ENV SMTP_SERVER=localhost
ENV SMTP_PORT=1025
ENV TO_MAILBOX=no-reply@localhost
ENV HTTP_PORT=8001

# Start the server app by default
CMD ./feedback-http-server \
      -smtp-server $SMTP_SERVER \
      -smtp-port $SMTP_PORT \
      -to $TO_MAILBOX \
      -http-port $HTTP_PORT

# Expose the server port
EXPOSE $HTTP_PORT
