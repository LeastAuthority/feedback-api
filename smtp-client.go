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
	"net/mail"
	"crypto/tls"
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

func connectAndSendEmail(hostname string, port uint, fromAddr string, toAddr string, subject string, body []byte) {
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
	auth := smtp.PlainAuth("", username, password, hostname)

	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
			ServerName: hostname,
		}
	conn, err := tls.Dial("tcp", hostPortStr, tlsconfig)

	c, err := smtp.NewClient(conn, hostname)
	if err != nil {
		log.Panic(err)
	}

	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	from := mail.Address{"", fromAddr}
	to   := mail.Address{"", toAddr}

	headers := make(map[string]string)

	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + emailBody

	log.Printf("sending email via %s to %s\n", hostPortStr, to)
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

	log.Printf("sent")
}
