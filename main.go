package main

// Listens on an HTTP port for POST method.
// if the post JSON payload is valid, it is emailed to a fixed address.
// Address to email to is configurable.

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	smtpHost  string
	smtpPort  uint
	to        string
	from      string
	subject   string
}

func (c *Config) handlePost(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling a post request to feedback url")
	fmt.Fprintf(w, "post\n")

	// take req.Body and pass it through a JSON decoder and turn
	// it into a feedback value.
	var body []byte
	_, err := req.Body.Read(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	req.Body.Close()
	go connectAndSendEmail(c.smtpHost, c.smtpPort, c.from, c.to, c.subject, string(body))
}

func main() {
	toAddressPtr := flag.String("to", "feedback@winden.app", "email address to which feedback is to be sent")
	smtpRelayHost := flag.String("smtp-server", "smtp.gmail.com", "smtp server that routes the email")
	smtpRelayPort := flag.Uint("smtp-port", 465, "smtp server port number")
	flag.Parse()

	c := Config{
		to: *toAddressPtr,
		from: "do-not-reply@winden.app",
		subject: "Winden Feedback",
		smtpPort: *smtpRelayPort,
		smtpHost: *smtpRelayHost,
	}
	// XXX: parse the email address to make sure it is a valid one.
	log.Printf("feedback email would be send to the address: %s\n", *toAddressPtr)

	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.handlePost).Methods("POST")


	srv := &http.Server{
		Addr: ":8001",
		Handler: r,
	}

	srv.ListenAndServe()
}
