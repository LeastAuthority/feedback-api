# Privacy friendly feedback API

## Why

We want to collect feedback from Users in more privacy friendly way, without collecting additional ourselves data or denying to do for 3rd parties.

## What

Simple HTTPS API backend service, which will get data in JSON format and send it to internal email (v1) to read.

## Primary Use cases

1. Allow users of Winden.app to leave extensive feedback
2. Allow users of Winden.app to rate and leave small feedback after each successful transfer by sender or receiver

## Usage

 - Build the code

 `go build`

 - Start the server with the address to which emails need to be sent

 `./feedback-api -to "foo@bar.org"`

 Server listens on `localhost:8001`.

 - Issue Post request:

 `curl --request POST --header "Content-Type: application/json" --data '{"q1": "a1", "q2": "a2"}' localhost:8001/v1/feedback`

## Docker image build

- Build image
`ocker build -t feedback . --build-arg var_SMTP_SERVER=localhost --build-arg var_SMTP_PORT=1025 --build-arg var_TO_MAILBOX=no-reply@localhost`

- Run image
`docker run -p 8001:8001 -t feedback-api`

## Technical aspect

- Use JSON format [RFC 8259](https://www.rfc-editor.org/rfc/rfc8259.html)

