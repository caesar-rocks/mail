package mailer

import (
	"crypto/tls"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type SMTPCfg struct {
	Username  string
	Password  string
	Host      string
	Port      string
	KeepAlive bool
	Timeout   int
	UseTLS    bool

	// By default, the encryption is set to SSL/TLS.
	Encryption *mail.Encryption
}

type smtpMailer struct {
	smtpClient *mail.SMTPClient
}

func newSMTP(cfg SMTPCfg) (MailerClient, error) {
	server := mail.NewSMTPClient()

	server.Host = cfg.Host
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return nil, err
	}
	server.Port = port
	server.Username = cfg.Username
	server.Password = cfg.Password
	if cfg.Encryption != nil {
		server.Encryption = *cfg.Encryption
	} else {
		server.Encryption = mail.EncryptionSSLTLS
	}
	server.KeepAlive = cfg.KeepAlive
	server.ConnectTimeout = time.Duration(cfg.Timeout) * time.Second
	server.SendTimeout = time.Duration(cfg.Timeout) * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.UseTLS}

	smtpClient, err := server.Connect()

	if err != nil {
		return nil, err
	}

	return &smtpMailer{smtpClient: smtpClient}, nil
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
