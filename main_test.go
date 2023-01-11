package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

func TestFailedServerConnection(t *testing.T) {
	t.Skip()
}

func TestTooLargeDataSent(t *testing.T) {
	t.Skip()
}

func TestNonJsonRejection(t *testing.T) {
	t.Skip()
}

func TestSeveralFeedbacksSending(t *testing.T) {
	t.Skip()
}

func TestWrongUrlResponse(t *testing.T) {
	t.Skip()
}

func TestSMTPServerReconnection(t *testing.T) {
	t.Skip()
}

func TestSimpleSendEmail(t *testing.T) {
	t.Setenv("SMTP_USE_TLS", "false")
	t.Setenv("SMTP_USERNAME", "test")

	//Mock smtp server
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
	})

	// To start server use Start() method
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}

	// Server's port will be assigned dynamically after server.Start()
	// for case when portNumber wasn't specified
	hostAddress, portNumber := "127.0.0.1", server.PortNumber()

	req, err := http.NewRequest("POST", "/v1/feedback", bytes.NewBuffer([]byte(`{"feedback": {"questions": [{"question": "What's great (if anything)?","answer": "I like speed."}]}}`)))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	c := Config{
		to:       "feedback@test.test",
		from:     "no-reply@test.test",
		subject:  "Feedback",
		smtpPort: uint(portNumber),
		smtpHost: hostAddress,
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	msg := server.Messages()[0].MsgRequest()
	expectMsg := "Q: What's great (if anything)?"

	if !strings.Contains(msg, expectMsg) {
		t.Errorf("Invalid message, got:\n %v, \nexpected:\n %v",
			msg, expectMsg)
	}

	if err := server.Stop(); err != nil {
		fmt.Println(err)
	}
}
