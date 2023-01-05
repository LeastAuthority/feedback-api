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

 - Set two environment variables while invoking the `feedback-api` executable, `SMTP_USERNAME` and `SMTP_PASSWORD`.
 - Start the server with the address to which emails need to be sent

 `./feedback-api -to "foo@bar.org"`

 or

 `SMTP_USERNAME="foo@foobar.in" SMTP_PASSWORD="barbazquux" ./feedback-api -to foo@barbaz.com -smtp-server smtp.abc.xyz -smtp-port 465`

If you are using `bash`, before typing in the command above, type a `SPC` character, so that the above command carrying the username and password won't get into the bash history.

 Server listens on `localhost:8001`.

 - Issue Post request:

 `curl --request POST --header "Content-Type: application/json" --data '{"q1": "a1", "q2": "a2"}' localhost:8001/v1/feedback`

### Disable TLS

You can disable TLS by setting the environment variable `SMTP_USE_TLS` to `"false"`. This can be useful if you want to use a dummy SMTP server for local development.

## Technical aspect

- Use JSON format [RFC 8259](https://www.rfc-editor.org/rfc/rfc8259.html)

