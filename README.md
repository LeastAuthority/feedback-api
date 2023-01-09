# Privacy-friendly feedback API

## Why

We want to collect feedback from Users in a more privacy-friendly way, without collecting additional data or denying to do for 3rd parties.

## What

Simple HTTPS API backend service, which will get data in JSON format and send it to an internal email (v1) to read.

## Primary Use cases

1. Allow users of the Winden.app to leave extensive feedback
2. Allow users of the Winden.app to rate and leave small feedback after each successful transfer by the sender or receiver

## Usage

 - Build the code

 `go build`

 - Set two environment variables while invoking the `feedback-api` executable, `SMTP_USERNAME` and `SMTP_PASSWORD`.
 - Start the server with the address to which emails need to be sent

 `./feedback-api -to "foo@bar.org"`

 or

 `SMTP_USERNAME="foo@foobar.in" SMTP_PASSWORD="barbazquux" ./feedback-api -to foo@barbaz.com -smtp-server smtp.abc.xyz -smtp-port 465`

If you are using `bash`, before typing in the command above, type a `SPC` character, so that the above command carrying the username and password won't get into the bash history.
A server listens on `localhost:8001`.

 - Issue Post request:

 `curl --request POST --header "Content-Type: application/json" --data '{"q1": "a1", "q2": "a2"}' localhost:8001/v1/feedback`

### Disable TLS

You can disable TLS by setting the environment variable `SMTP_USE_TLS` to `"false"`. This can be useful if you want to use a dummy SMTP server for local development.

## Docker image build

- Build image
`docker build -t feedback-api .`

- Run image with default localhost SMTP (note: SMTP server should run inside it)
`docker run -p 8001:8001 -t feedback-api`

- Run image with overwritte SMTP config
`docker run -p 8001:8001 -e SMTP_SERVER=smtp.example.com -e SMTP_PORT=25 -e TO_MAILBOX=feedback@example.com -t feedback-api`

- Run local SMTP server for development and testing (maildev)[https://github.com/maildev/maildev]
`docker run -p 1080:1080 -p 1025:1025 maildev/maildev`  

- Run image and connect to local SMTP server (maildev)[https://github.com/maildev/maildev]
`docker run -p 8001:8001 -e SMTP_SERVER=<<workstation_ip>> -t feedback-api`  

## Technical aspect

- Use JSON format [RFC 8259](https://www.rfc-editor.org/rfc/rfc8259.html)

