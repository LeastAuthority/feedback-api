package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"encoding/json"
	"text/template"
	"bytes"
	"net/smtp"
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

func connectAndSendEmail(hostname string, port uint, from string, to string, subject string, password string, body []byte) {
	emailBody, err := parseBody(body)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	if username == "" {
		log.Printf("WARNING: SMTP_USERNAME not set")
		return
	}
	hostPortStr := fmt.Sprintf("%s:%s", hostname, strconv.Itoa(int(port)))
	auth := smtp.PlainAuth("", from, password, hostPortStr)

	msg :=  "From: " + from + " \n"+
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		emailBody + "\r\n"

	log.Printf("sending email to %s\n", to)
	err = smtp.SendMail(hostPortStr,
		auth,
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Printf("sent")
}
