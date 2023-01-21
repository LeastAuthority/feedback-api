# Privacy-friendly feedback API

[![Report](https://goreportcard.com/report/github.com/LeastAuthority/feedback-api?style=flat)](https://goreportcard.com/report/github.com/LeastAuthority/feedback-api)

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

 `curl --request POST --header "Content-Type: application/json" --data '{"feedback":{"questions":[{"question":"q1","answer":"a1"},{"question":"q2","answer":"a2"}]}}' localhost:8001/v1/feedback`

### Disable TLS

You can disable TLS by setting the environment variable `SMTP_USE_TLS` to `"false"`. This can be useful if you want to use a dummy SMTP server for local development.

## Docker

To facilitate integration and testing, Docker Composer support is provided in this repo.

### Requirements

- Docker (20.10.5+)
- Docker Composer (1.25.0+)

### Configuration

The relevant environment variables described above can be adapted a local `.env` file:

```
SMTP_SERVER=smtp-server1.local
SMTP_PORT=1025
SMTP_USERNAME=xxx
SMTP_PASSWORD=xxx
SMTP_USE_TLS=false
TO_MAILBOX=no-reply@localhost
HTTP_PORT=8001
```

### Usage

It is recommended to explicitly build the images at least once, then subsequently in case of changes on Docker config and/or images:

```
docker-compose [--compatibility] build
```


The following command should start the required services:

```
docker-compose [--compatibility] up
```

REM: Use the `--compatibility` switch (if supported) whenever resources need to be limitted (CPU and MEM).

From this point, the service(s) (http and smtp) should be running and ready to be tested ( see `curl` command above).

### Advanced Usage

If needed, Docker compose can also be used to start a single service. For instance, here is how to run a standalone HTTP server w/o the SMTP service:

```
docker-compose [--compatibility] run \
--no-deps \
--rm \
--publish 8001:8001 \
-e SMTP_SERVER=smtp.example.com \
-e SMTP_PORT=465 \
-e SMTP_USE_TLS=true \
http-server1
```

## Technical aspect

- Use JSON format [RFC 8259](https://www.rfc-editor.org/rfc/rfc8259.html)
