package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

func TestParseBody(t *testing.T) {

	exampleFull, err := os.ReadFile("examples/full-feedback.json")
	if err != nil {
		t.Errorf("Cannot read example file")
	}
	testCases := []struct {
		name        string
		input       []byte
		expected    string
		expectError bool
	}{
		{
			name:  "valid json",
			input: exampleFull,
			expected: `Title: Full Feedback Form

From: test.app

Rate: 10 (numbers)

Q: What's great (if anything)?
A: I like speed.

Q: What do you find product useful for?
A: To transfer personal files.

Q: What's missing or what's not great?
A: Ability to do multiple file transfer.
`,
		},
		{
			name:  "valid no rating",
			input: []byte(`{"channel":"test.app","feedback":{"title":"Full Feedback Form","rate":{"type":"numbers","value":null},"questions":[{"question":"What's great (if anything)?","answer":"Nothing"},{"question":"What do you find product useful for?","answer":""},{"question":"What's missing or what's not great?","answer":""}]}}`),
			expected: `Title: Full Feedback Form

From: test.app

Rate: not rated

Q: What's great (if anything)?
A: Nothing

Q: What do you find product useful for?
A: 

Q: What's missing or what's not great?
A: 
`,
		},
		{
			name:        "invalid json",
			input:       []byte(`{"Questions": [`),
			expectError: true,
		},
		{
			name:        "json, incorrect template",
			input:       []byte(`{"Test": []}`),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseBody(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected:\n%q\n, but got:\n%q\n", tc.expected, result)
			}
		})
	}
}

func TestConnectAndSendEmailTls(t *testing.T) {
	t.Setenv("SMTP_USE_TLS", "true")
	t.Setenv("SMTP_USERNAME", "test")

	// Write Tls verification test
	hostname := "smtp.gmail.com"
	port := uint(465)
	fromAddr := "feedback@test.test"
	toAddr := "no-reply@test.test"
	subject := "Test Email"
	body := payload

	err := connectAndSendEmail(hostname, port, fromAddr, toAddr, subject, body)

	expectMsg := "Username and Password not accepted"

	if !strings.Contains(err.Error(), expectMsg) {
		t.Errorf("Invalid message, got:\n %v, \nexpected:\n %v",
			err, expectMsg)
	}
}
func TestConnectAndSendEmailInsecureTls(t *testing.T) {
	t.Setenv("SMTP_USE_INSECURE_TLS", "true")
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

	hostAddress, portNumber := "127.0.0.1", server.PortNumber()

	// Write Tls verification test
	hostname := hostAddress
	port := uint(portNumber)
	fromAddr := "feedback@test.test"
	toAddr := "no-reply@test.test"
	subject := "Test Email"
	body := payload

	err := connectAndSendEmail(hostname, port, fromAddr, toAddr, subject, body)

	// Not the best verification, however library doesn't have yet TLS support
	expectMsg := "first record does not look like a TLS handshake"

	if !strings.Contains(err.Error(), expectMsg) {
		t.Errorf("Invalid message, got:\n %v, \nexpected:\n %v",
			err, expectMsg)
	}

	if err := server.Stop(); err != nil {
		fmt.Println(err)
	}
}
