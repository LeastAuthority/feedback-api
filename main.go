package main

// Listens on an HTTP port for POST method.
// if the post JSON payload is valid, it is emailed to a fixed address.
// Address to email to is configurable.

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"io/ioutil"

	"github.com/gorilla/mux"
)

type Config struct {
	smtpHost  string
	smtpPort  uint
	to        string
	from      string
	subject   string
	password  string
}

const (
	MaxPayloadSize = 32768
)

func (c *Config) sendEmail(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling a post request to feedback url")

	// take req.Body and pass it through a JSON decoder and turn
	// it into a feedback value.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading the request body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	if len(body) > MaxPayloadSize {
		// XXX What should be the HTTP error here? For sure
		// 4xx since this is an error on Client's part. Now,
		// 400 or 413? Will go with 400 for now, but this is
		// something to revisit..
		log.Printf("payload size is larger than 32kB\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check for "syntax" errors but do not decode into Go
	// values
	if !json.Valid(body) {
		log.Printf("malformed JSON payload\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go connectAndSendEmail(c.smtpHost, c.smtpPort, c.from, c.to, c.subject, c.password, body)
}

func main() {
	toAddressPtr := flag.String("to", "feedback@winden.app", "email address to which feedback is to be sent")
	smtpRelayHost := flag.String("smtp-server", "smtp.gmail.com", "smtp server that routes the email")
	smtpRelayPort := flag.Uint("smtp-port", 465, "smtp server port number")
	smtpPassword := flag.String("smtp-server-password", "", "smtp server password")
	flag.Parse()

	c := Config{
		to: *toAddressPtr,
		from: "doNotReply@leastauthority.com",
		subject: "Winden Feedback",
		smtpPort: *smtpRelayPort,
		smtpHost: *smtpRelayHost,
		password: *smtpPassword,
	}
	// XXX: parse the email address to make sure it is a valid one.
	log.Printf("feedback email would be send to the address: %s\n", *toAddressPtr)

	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")


	srv := &http.Server{
		Addr: ":8001",
		Handler: r,
	}

	srv.ListenAndServe()
}
