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
	var fullFeedback Feedback
	err := json.Unmarshal(body, &fullFeedback)
	if err != nil {
		log.Printf("json parsing of the body failed\n")
		return "", err
	}
	feedbackTmpl := `
{{- range .Questions}}
Q: {{.Question}}
A: {{.Answer}}
{{end}}
`
	output := bytes.NewBufferString("")
	tmpl := template.Must(template.New("full feedback template").Parse(feedbackTmpl))
	tmpl.Execute(output, &fullFeedback.Full)

	log.Printf("DEBUG: %s\n", output.String())

	return output.String(), nil
}

func connectAndSendEmail(hostname string, port uint, from string, to string, subject string, body []byte) {
	emailBody, err := parseBody(body)
	if err != nil {
		log.Printf("%v\n", err)
	}

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
		"%s\r\n", to, subject, from, emailBody)
	email := strings.NewReader(msg)
	log.Printf("sending email to %s\n", to)
	err = smtpClient.SendMail(from, []string{to}, email)

	if err != nil {
		log.Printf("%v\n", err)
	}
	smtpClient.Close()

}
