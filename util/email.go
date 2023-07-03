package util

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"time"
)

func SendEmail(from, pwd, to, subject, body string) error {
	e := &email.Email{
		To:      []string{to},
		From:    fmt.Sprintf("Solvery <%s>", from),
		Subject: subject,
		Text:    []byte(body),
		Headers: textproto.MIMEHeader{},
	}

	p, err := email.NewPool(
		"smtp.gmail.com:587",
		1,
		smtp.PlainAuth("", from, pwd, "smtp.gmail.com"),
	)

	if err != nil {
		return err
	}

	return p.Send(e, time.Second*5)
}
