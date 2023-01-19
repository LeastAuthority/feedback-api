package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"text/template"
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
	err = tmpl.Execute(output, &fullFeedback.Full)
	res := output.String()
	if err != nil {
		log.Panic(err)
		return "", err
	}

	if len(res) <= 1 {
		err := errors.New("json doesn't match template")
		return "", err
	}

	return res, nil
}

func connectAndSendEmail(hostname string, port uint, fromAddr string, toAddr string, subject string, body []byte) error {

	emailBody, err := parseBody(body)
	if err != nil {
		return err
	}

	useTls, err := strconv.ParseBool(os.Getenv("SMTP_USE_TLS"))
	if err != nil {
		useTls = true
	}
	useInsecureTls, err := strconv.ParseBool(os.Getenv("SMTP_USE_INSECURE_TLS"))
	if err != nil {
		useInsecureTls = false
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	if username == "" {
		log.Printf("WARNING: SMTP_USERNAME not set")
	}
	hostPortStr := fmt.Sprintf("%s:%s", hostname, strconv.Itoa(int(port)))
	auth := smtp.PlainAuth("", username, password, hostname)

	var conn net.Conn
	if useTls || useInsecureTls {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: useInsecureTls,
			ServerName:         hostname,
		}
		conn, err = tls.Dial("tcp", hostPortStr, tlsconfig)
		if err != nil {
			return err
		}
	} else {
		conn, err = net.Dial("tcp", hostPortStr)
		if err != nil {
			return err
		}
	}

	c, err := smtp.NewClient(conn, hostname)
	if err != nil {
		log.Panic(err)
	}

	if useTls {
		err = c.Auth(auth)
		if err != nil {
			log.Panic(err)
		}
	}

	from := mail.Address{Name: "Feedback Server", Address: fromAddr}
	to := mail.Address{Name: "Feedback Mailbox", Address: toAddr}

	headers := make(map[string]string)

	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
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

	err = c.Quit()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("sent")
	return nil
}
