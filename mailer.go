package mailer

import (
	"sync"

	"github.com/caesar-rocks/core"
)

type APIServiceType string

const (
	// mail version
	version                   = "0.1.0"
	SMTP       APIServiceType = "smtp"
	SENDGRID   APIServiceType = "sendgrid"
	MAILGUN    APIServiceType = "mailgun"
	AMAZON_SES APIServiceType = "amazon-ses"
	RESEND     APIServiceType = "resend"
)

var (
	mailer *Mail
	once   sync.Once
)

// Mail is a struct that holds the configuration for the mailer.
type Mail struct {
	host         string
	port         string
	username     string
	password     string
	apiService   APIServiceType
	apiKey       string
	emailToSend  chan MailerMessage
	mailErr      chan error
	keepAlive    bool
	timeout      int
	mailerClient Mailer
}

// MailConfig is a struct that holds the configuration for the mailer.
type MailConfig struct {
	// FromName is the name that will be used as the sender.
	FromName string `json:"from_name,omitempty"`
	// ReplyToEmail is the name that will be used as the sender.
	ReplyToEmail string `json:"reply_to_email,omitempty"`
	// Host is the host of the mail server.
	Host string `json:"host,omitempty"`
	// HostUser is the username for the mail server.
	HostUser string `json:"host_user,omitempty"`
	// HostPassword is the password for the mail server.
	HostPassword string `json:"host_password,omitempty"`
	// Port is the port of the mail server.
	Port string `json:"port,omitempty"`
	// UseTLS is a boolean that determines whether to use TLS.
	UseTLS bool `json:"use_tls,omitempty"`
	// UseSSL is a boolean that determines whether to use SSL.
	UseSSL bool `json:"use_ssl,omitempty"`
	// Timeout is the timeout to connect to SMTP Server and to send the email and wait respond
	Timeout int `json:"timeout,omitempty"`
	// APIService is the service to use for sending emails.
	APIService APIServiceType `json:"api_service,omitempty"`
	// APIKey is the key to use for sending emails.
	APIKey string
	// KeepAlive to keep alive connection
	KeepAlive bool
	// MailerClient is the mailer client to use for sending emails.
	mailerClient Mailer
}

// NewMailer creates a new mailer instance. It is a singleton.
// It requires a MailConfig struct as an argument.
func NewMailer(opt MailConfig) *Mail {
	core.ValidateEnvironmentVariables[MailConfigEnv]()
	once.Do(func() {
		mailer = &Mail{
			host:         opt.Host,
			port:         opt.Port,
			username:     opt.HostUser,
			password:     opt.HostPassword,
			apiService:   opt.APIService,
			apiKey:       opt.APIKey,
			keepAlive:    opt.KeepAlive,
			timeout:      opt.Timeout,
			emailToSend:  make(chan MailerMessage, 200),
			mailErr:      make(chan error),
			mailerClient: getMailerClient(opt),
		}
		go mailer.listenForEmails()
	})
	return mailer
}

// Send sends an email message using the chosen API service.
func (m *Mail) Send(message MailerMessage) error {
	mailer.emailToSend <- message
	return <-mailer.mailErr
}

// Close closes the emailToSend, result channels and the mailerClient.
func (m *Mail) Close() {
	close(m.emailToSend)
	close(m.mailErr)
	m.mailerClient.Close()
}

// send sends the email message using the chosen API service.
func (m *Mail) send(msg MailerMessage) error {
	return m.mailerClient.Send(msg)
}

// ListenForEmails listens for email messages and sends them using the chosen API service.
// It is a blocking function that should be run in a goroutine.
func (m *Mail) listenForEmails() {
	for {
		select {
		case msg, ok := <-m.emailToSend:
			if !ok {
				return
			}
			err := m.send(msg)
			m.mailErr <- err
		}
	}
}
