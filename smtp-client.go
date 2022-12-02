package main

import (
	"strconv"
	"fmt"
	"log"
	"strings"

	"github.com/emersion/go-smtp"
)

func connectAndSendEmail(hostname string, port uint, from string, to string, subject string, body string) {
	hostPortStr := fmt.Sprintf("%s:%s", hostname, strconv.Itoa(int(port)))

	smtpClient, err := smtp.Dial(hostPortStr)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer smtpClient.Close()

	msg := fmt.Sprintf("To: %s\r\n" +
		"Subject: %s\r\n" +
		"From: %s\r\n" +
		"\r\n" +
		"%s\r\n", to, subject, from, body)
	email := strings.NewReader(msg)
	log.Printf("sending email to %s\n", to)
	err = smtpClient.SendMail(from, []string{to}, email)

	if err != nil {
		log.Printf("%v\n", err)
	}
}
