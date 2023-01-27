package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

var payload = []byte(`{"channel":"test.app","feedback": {"title": "feedback","rate": {"type" : "numbers","value" : 10},"questions": [{"question": "What's great (if anything)?","answer": "I like speed."}]}}`)

var payloadDynamic = (`{"channel":"test.app","feedback": {"title": "feedback","rate": {"type" : "numbers","value" : 10},"questions": [{"question": "%d","answer": "%d"}]}}`)

func TestFailedServerConnection(t *testing.T) {

	req, err := http.NewRequest("POST", "/v1/feedback", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	c := Config{
		to:       "feedback@test.test",
		from:     "no-reply@test.test",
		subject:  "Feedback",
		smtpPort: 12345,
		smtpHost: "localhost",
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}

}

func TestTooLargeDataSent(t *testing.T) {

	payloadLarge := make([]byte, 40000)

	_, err := rand.Read(payloadLarge)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/feedback", bytes.NewBuffer(payloadLarge))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	c := Config{
		to:       "feedback@test.test",
		from:     "no-reply@test.test",
		subject:  "Feedback",
		smtpPort: 12345,
		smtpHost: "localhost",
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

func TestNonJsonRejection(t *testing.T) {
	payloadNonJson := make([]byte, 1024)
	_, err := rand.Read(payloadNonJson)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/feedback", bytes.NewBuffer(payloadNonJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	c := Config{
		to:       "feedback@test.test",
		from:     "no-reply@test.test",
		subject:  "Feedback",
		smtpPort: 12345,
		smtpHost: "localhost",
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

func TestSeveralFeedbacksSending(t *testing.T) {
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

	// send 3 messages
	for i := range [3]int{} {

		_payload := fmt.Sprintf(payloadDynamic, i, i)

		req, err := http.NewRequest("POST", "/v1/feedback", bytes.NewBuffer([]byte(_payload)))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		if err != nil {
			t.Fatal(err)
		}
		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		msg := server.Messages()[i].MsgRequest()
		expectMsg := fmt.Sprintf(`Q: %d`, i)

		if !strings.Contains(msg, expectMsg) {
			t.Errorf("Invalid message, got:\n %v, \nexpected:\n %v",
				msg, expectMsg)
		}
	}

	if err := server.Stop(); err != nil {
		fmt.Println(err)
	}
}

func TestWrongUrlRequest(t *testing.T) {
	payloadRandom := make([]byte, 1024)
	_, err := rand.Read(payloadRandom)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payloadRandom))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	c := Config{
		to:       "feedback@test.test",
		from:     "no-reply@test.test",
		subject:  "Feedback",
		smtpPort: 12345,
		smtpHost: "localhost",
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/feedback", c.sendEmail).Methods("POST")

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusNotFound)
	}
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

	req, err := http.NewRequest("POST", "/v1/feedback", bytes.NewBuffer(payload))
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
