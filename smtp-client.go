package main

import (
	"strconv"
	"fmt"
	"strings"

	"github.com/emersion/go-smtp"
)

func connectAndSendEmail(hostname string, port uint, from string, to string, subject string, body string) error {
	hostPortStr := fmt.Sprintf("%s:%s", hostname, strconv.Itoa(int(port)))

	c, err := smtp.Dial(hostPortStr)
	if err != nil {
		return err
	}
	defer c.Close()

	msg := fmt.Sprintf("To: %s\r\n" +
		"Subject: %s\r\n" +
		"\r\n" +
		"%s\r\n", to, subject, body)
	email := strings.NewReader(msg)
	err = c.SendMail(from, []string{to}, email)

	return err
}
