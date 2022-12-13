package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"encoding/json"
	"text/template"
	"bytes"

	"github.com/emersion/go-smtp"
)

func parseBody(body []byte) (string, error) {
	var fullFeedback feedback
	err := json.Unmarshal(body, &fullFeedback)
	if err != nil {
		return "", err
	}

	feedbackTmpl := `
feedback:
{{- range $qanda := }}
`
	output := bytes.NewBufferString("")
	tmpl := template.Must(template.New("full feedback template").Parse(feedbackTmpl))
	tmpl.Execute(output, &fullFeedback.full.questions)

	return output.String(), nil
}

func connectAndSendEmail(hostname string, port uint, from string, to string, subject string, body []byte) {
	hostPortStr := fmt.Sprintf("%s:%s", hostname, strconv.Itoa(int(port)))

	smtpClient, err := smtp.Dial(hostPortStr)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"From: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, from, body)
	email := strings.NewReader(msg)
	log.Printf("sending email to %s\n", to)
	err = smtpClient.SendMail(from, []string{to}, email)

	if err != nil {
		log.Printf("%v\n", err)
	}
	smtpClient.Close()

}
