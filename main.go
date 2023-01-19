package main

// Listens on an HTTP port for POST method.
// if the post JSON payload is valid, it is emailed to a fixed address.
// Address to email to is configurable.

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
)

type Config struct {
	smtpHost    string
	smtpPort    uint
	to          string
	from        string
	subject     string
	httpPort    uint
	insecureTls bool
}

const (
	MaxPayloadSize = 32768
)

func (c *Config) sendEmail(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling a post request to feedback url")

	// take req.Body and pass it through a JSON decoder and turn
	// it into a feedback value.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading the request body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	if len(body) > MaxPayloadSize {
		// 413 for too large payload
		log.Printf("payload size is larger than 32kB\n")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	// check for "syntax" errors but do not decode into Go
	// values
	if !json.Valid(body) {
		log.Printf("malformed JSON payload\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respCh := make(chan error)
	go func() {
		err := connectAndSendEmail(c.smtpHost, c.smtpPort, c.from, c.to, c.subject, body)

		if err != nil {
			log.Printf("Failed sending feedback, error: %s\n", err)
			respCh <- err
		}
		respCh <- nil
	}()

	// in case failed to send, respond with error
	if err = <-respCh; err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	toAddressPtr := flag.String("to", "feedback@localhost", "email address to which feedback is to be sent")
	smtpRelayHost := flag.String("smtp-server", "localhost", "smtp server that routes the email")
	smtpRelayPort := flag.Uint("smtp-port", 1025, "smtp server port number")
	httpPort := flag.Uint("http-port", 8001, "HTTP server port number")
	flag.Parse()

	c := Config{
		to:       *toAddressPtr,
		from:     "no-reply@localhost",
		subject:  "Feedback",
		smtpPort: *smtpRelayPort,
		smtpHost: *smtpRelayHost,
		httpPort: *httpPort,
	}

	// email address validation
	_, err := mail.ParseAddress(*toAddressPtr)
	if err != nil {
		log.Println("Invalid email address")
		panic(err)
	}
	log.Printf("Feedback email will be sent to: %s\n", *toAddressPtr)

	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.httpPort),
		Handler: r,
	}

	err = srv.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
