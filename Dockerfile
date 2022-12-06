FROM golang:1.19.3-alpine

COPY . /app
WORKDIR /app

RUN go build

ARG required var_SMTP_SERVER
ENV SMTP_SERVER=${var_SMTP_SERVER}

ARG required var_SMTP_PORT
ENV SMTP_PORT=${var_SMTP_PORT}

ARG required var_TO_MAILBOX
ENV TO_MAILBOX=${var_TO_MAILBOX}

CMD ./feedback-api -smtp-server $SMTP_SERVER -smtp-port $SMTP_PORT -to $TO_MAILBOX
EXPOSE 8001
