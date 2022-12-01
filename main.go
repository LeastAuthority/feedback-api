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

	log.Printf("handling a post request to feedback url")
	fmt.Fprintf(w, "post\n")

	// take req.Body and pass it through a JSON decoder and turn
	// it into a feedback value.
}

func main() {
	toAddressPtr := flag.String("to", "feedback@winden.app", "email address to which feedback is to be sent")
	flag.Parse()

	// XXX: parse the email address to make sure it is a valid one.
	log.Printf("feedback email would be send to the address: %s\n", *toAddressPtr)

	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", handlePost).Methods("POST")


	srv := &http.Server{
		Addr: ":8001",
		Handler: r,
	}

	srv.ListenAndServe()
}
