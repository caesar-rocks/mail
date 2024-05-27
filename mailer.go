package mailer

type APIServiceType string

// MailConfig is a struct that holds the configuration for the mailer.
type MailConfig struct {
	// FromName is the name that will be used as the sender.
	FromName string
	// ReplyToEmail is the name that will be used as the sender.
	ReplyToEmail string
	// Host is the host of the mail server.
	Host string
	// HostUser is the username for the mail server.
	HostUser string
	// HostPassword is the password for the mail server.
	HostPassword string
	// Port is the port of the mail server.
	Port string
	// UseTLS is a boolean that determines whether to use TLS.
	UseTLS bool
	// UseSSL is a boolean that determines whether to use SSL.
	UseSSL bool
	// Timeout is the timeout to connect to SMTP Server and to send the email and wait respond
	Timeout int
	// APIService is the service to use for sending emails.
	APIService APIServiceType
	// APIKey is the key to use for sending emails.
	APIKey string
	// KeepAlive to keep alive connection
	KeepAlive bool
	// MailerClient is the mailer client to use for sending emails.
	mailerClient Mailer
}

// Mail is a struct that holds the mailer instance and its configurations
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

// NewMailer creates a new mailer instance.
func NewMailer(cfg MailConfig) *Mail {
	mailer := &Mail{
		host:         cfg.Host,
		port:         cfg.Port,
		username:     cfg.HostUser,
		password:     cfg.HostPassword,
		apiService:   cfg.APIService,
		apiKey:       cfg.APIKey,
		keepAlive:    cfg.KeepAlive,
		timeout:      cfg.Timeout,
		emailToSend:  make(chan MailerMessage, 200),
		mailErr:      make(chan error),
		mailerClient: getMailerClient(cfg),
	}

	go mailer.listenForEmails()

	return mailer
}

// Send sends an email message using the chosen API service.
func (m *Mail) Send(msg MailerMessage) error {
	m.emailToSend <- msg
	return <-m.mailErr
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
