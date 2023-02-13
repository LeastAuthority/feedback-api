# Start from the official Go image based on Alpine
FROM golang:1.19.3-alpine

LABEL Description="HTTP server for feedback API"

# Parameters for default user:group
ARG uid=1000
ARG user=appuser
ARG gid=1000
ARG group=appgroup

# Add user and group for build and runtime
RUN id ${user} > /dev/null 2>&1 || \
    { addgroup -g "${gid}" "${group}" && adduser -D -h /home/${user} -s /bin/bash -G "${group}" -u "${uid}" "${user}"; }

# Prepare directories
RUN DIRS="/src /app" && \
    mkdir -p ${DIRS} && \
    chown -R ${user}:${group} $DIRS

# Switch to non-root user
USER ${user}

# Inject the source code
COPY *.go go.* /src/

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
ENV SMTP_FROM=no-reply@localhost
ENV TO_MAILBOX=feedback@localhost
ENV HTTP_PORT=8001

# Start the server app by default
CMD ./feedback-http-server \
      -smtp-server $SMTP_SERVER \
      -smtp-port $SMTP_PORT \
      -from $SMTP_FROM \
      -to $TO_MAILBOX \
      -http-port $HTTP_PORT

# Expose the server port
EXPOSE $HTTP_PORT
