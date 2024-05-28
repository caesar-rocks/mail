package mailer

import (
	"crypto/tls"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type smtpParams struct {
	Username  string
	Password  string
	Host      string
	Port      string
	KeepAlive bool
	Timeout   int
	useTLS    bool
}

type smtpMailer struct {
	smtpClient *mail.SMTPClient
}

func newSMTP(params smtpParams) MailerClient {
	server := mail.NewSMTPClient()

	server.Host = params.Host
	server.Port = getPort(params.Port)
	server.Username = params.Username
	server.Password = params.Password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = params.KeepAlive
	server.ConnectTimeout = time.Duration(params.Timeout) * time.Second
	server.SendTimeout = time.Duration(params.Timeout) * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: params.useTLS}

	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}
	return &smtpMailer{smtpClient: smtpClient}
}

func (m *smtpMailer) Send(msg Mail) error {
	_, email := buildMessage(msg)

	if email.Error != nil {
		return email.Error
	}

	err := email.Send(m.smtpClient)
	if err != nil {
		return err
	}
	return nil
}

func (m *smtpMailer) Close() {
	m.smtpClient.Close()
}
