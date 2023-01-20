# Start from the official Go image based on Alpine
FROM golang:1.19.3-alpine

# Inject the source code
COPY . /app
WORKDIR /app

# Download dependencies and build the server app
RUN go mod download && go mod verify
RUN go build -buildvcs=false -v

# Default smtp config
ENV SMTP_SERVER=localhost
ENV SMTP_PORT=1025
ENV TO_MAILBOX=no-reply@localhost
ENV HTTP_PORT=8001

# Start the server app by default
CMD ./feedback-api \
      -smtp-server $SMTP_SERVER \
      -smtp-port $SMTP_PORT \
      -to $TO_MAILBOX \
      -http-port $HTTP_PORT

# Expose the server port
EXPOSE $HTTP_PORT
