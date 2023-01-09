FROM golang:1.19.3-alpine

COPY . /app
WORKDIR /app

RUN go build

# Default smtp config
ENV SMTP_SERVER=localhost
ENV SMTP_PORT=1025
ENV TO_MAILBOX=no-reply@localhost
ENV HTTP_PORT=8001

CMD ./feedback-api \
      -smtp-server $SMTP_SERVER \
      -smtp-port $SMTP_PORT \
      -to $TO_MAILBOX \
      -http-port $HTTP_PORT

EXPOSE $HTTP_PORT
