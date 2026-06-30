package smtp

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTP struct {
	host string
	port string
	password string
	from string
}

func NewSMTP(
	host,
	port,
	password,
	from string,
) *SMTP {
	return  &SMTP{
		host: host,
		port: port,

		password: password,
		from: from,
	}
}

func (s *SMTP) SendOTP(
	ctx context.Context,
	email string,
	code string,
) error {
	auth := smtp.PlainAuth(
		"",
		s.from,
		s.password,
		s.host,
	)

	subject := "Subject: Email Verificartion\r\n"

	body := fmt.Sprintf(
		"Your verification code is: %s\r\n\n The code is valid for 5 minutes.",
		code,
	) 

	message := []byte(subject + "\r\n" + body)
	return  smtp.SendMail(
		s.host+":"+s.port,
		auth,
		s.from,
		[]string{email},
		message,
	)
}