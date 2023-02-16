# Privacy-friendly feedback API

[![Go Report Card](https://goreportcard.com/badge/github.com/LeastAuthority/feedback-api)](https://goreportcard.com/report/github.com/LeastAuthority/feedback-api)

## Why

We want to collect feedback from Users in a more privacy-friendly way, without collecting additional data or denying to do for 3rd parties.

## What

Simple HTTPS API backend service, which will get data in JSON format and send it to an internal email (v1) to read.

## Primary Use cases

1. Allow users of the Winden.app to leave extensive feedback
2. Allow users of the Winden.app to rate and leave small feedback after each successful transfer by the sender or receiver

## Usage

This section describe how to directly compile and run the application assuming the [Go](https://golang.org) programming language compiler is installed.

Alternatively the use of [Docker](#docker) Composer is described in the next section.

### Build

The application can be build as follows:

```
go build
```

### Run

The application directly compiled with Go can be started as follows:

```
./feedback-api \
-from "no-reply@example.com" \
-to foo@barbaz.com \
-smtp-server smtp.abc.xyz \
-smtp-port 465 \
-http-port 8001
```

As a result, an HTTP server should be listenning on `localhost:8001` (CTRL+C to stop).

### Configure

There are a few environment variables avaialble to configure the application:

- SMTP_HELO: To announce a valid client hostname when the SMTP client contacts the server.
- SMTP_USERNAME and SMTP_PASSWORD: To authenticate the SMTP client if the server supports it.
- SMTP_USE_TLS: Set to `false` to disable TLS, which can be useful if you want to use a dummy SMTP server for local development.
- SMTP_USE_INSECURE_TLS: Set to `true` to allow insecure TLS, which can be useful if you want to test with a self-sign certificate server in development env.

Example:

```
SMTP_HELO="gw.example.com" \
SMTP_USERNAME="foo@foobar.in" \
SMTP_PASSWORD="barbazquux" \
SMTP_USE_INSECURE_TLS=true \
./feedback-api \
-from "no-reply@example.com" \
-to foo@barbaz.com \
-smtp-server localhost
```

IMPORTANT: If you are using `bash`, before typing in the command above, type a `SPC` character, so that the above command carrying the username and password won't get into the bash history.


### Test

Once the application in running, a POST request can be sent:

 ```
 curl \
--request POST \
--header "Content-Type: application/json" \
--data '{"feedback":{"title": "Full Feedback Form","rate": {"type" : "numbers","value" : 10},"questions":[{"question":"q1","answer":"a1"},{"question":"q2","answer":"a2"}]}}' \
http://localhost:8001/v1/feedback
```

REM: The JSON content needs to match the defined [template](./templates.go)!

## Docker

To facilitate integration and testing, Docker Composer support is provided in this repo.

This also include a dummy SMTP server for testing purpose.

### Requirements

- Docker (20.10.5+)
- Docker Composer (1.25.0+)

### Configuration

This method relies exclusively on environment variables which means both the application parameters and the above described variables need to be defined in a local `.env` file (ignored by Git):

```
SMTP_FROM=no-reply@example.com
TO_MAILBOX=foo@barbaz.com
SMTP_SERVER=smtp.abc.xyz
SMTP_PORT=465
HTTP_PORT=8001

SMTP_HELO=http-server1.local
SMTP_USERNAME=foo@foobar.in
SMTP_PASSWORD=barbazquux
SMTP_USE_TLS=true
SMTP_USE_INSECURE_TLS=true
```

**REM**: Some variables can lead to conflict between containers (e.g.: `SMTP_SERVER`)! We may consider renaming them.

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

From this point, the service(s) (http and smtp) should be running and ready to be tested (see `curl` command above).

### Advanced Usage

If needed, Docker compose can also be used to start a single service.

For instance, here is how to run a standalone HTTP server w/o the SMTP service:

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
