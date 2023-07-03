package notification

import (
	"fmt"
	"net/smtp"
	"paniq-auth-user-service/pkg/contracts"
)

type Usecase struct {
	SmtpHost     string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
}

func (u *Usecase) NotifyWithEmail(from string, to string, subject string, message string) {

	auth := smtp.PlainAuth("", u.SmtpUser, u.SmtpPassword, u.SmtpHost)

	// Here we do it all: connect to our server, set up a message and send it
	tos := []string{to}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s \r\n\r\n%s", to, subject, message))

	smtp.SendMail(u.SmtpHost+":"+u.SmtpPort, auth, from, tos, msg)

	//var emailError chan error
	//go func(err chan error, host string, port string, auth smtp.Auth, from string, to []string, msg []byte) {
	//	err <- smtp.SendMail(host+":"+port, auth, from, to, msg)
	//}(emailError, u.SmtpHost, u.SmtpPort, auth, from, tos, msg)
}

func NewUsecase(
	SmtpHost string,
	SmtpPort string,
	SmtpUser string,
	SmtpPassword string,
) contracts.INotificationService {
	return &Usecase{
		SmtpHost:     SmtpHost,
		SmtpPort:     SmtpPort,
		SmtpUser:     SmtpUser,
		SmtpPassword: SmtpPassword,
	}
}
