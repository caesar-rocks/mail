package mailer

type APIServiceType string

const (
	SMTP       APIServiceType = "smtp"
	SENDGRID   APIServiceType = "sendgrid"
	MAILGUN    APIServiceType = "mailgun"
	AMAZON_SES APIServiceType = "amazon-ses"
	RESEND     APIServiceType = "resend"
)

type Attachment struct {
	// Name is the name of the attachment.
	Name string
	// Path is the path to the attachment.
	Path string
}

type Mail struct {
	// To is the email address of the recipient.
	To string
	// From is the email address of the sender. defaults to the host user.
	From string
	// Html is the html content of the email.
	Html string
	// Text is the text content of the email.
	Text string
	// Subject is the subject of the email.
	Subject string
	// Cc is the email address of the cc recipient.
	Cc string
	// Bcc is the email address of the bcc recipient.
	Bcc string
	// ReplyTo is the email address to reply to.
	ReplyTo string
	// Attachments is an array of attachments.
	Attachments []Attachment
}

type MailerClient interface {
	Send(msg Mail) error
	Close()
}

// MailCfg is a struct that holds the configuration for the mailer.
type MailCfg struct {
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
	// APISecret is the secret to use for sending emails.
	APISecret string
	// Region is the region to use for sending emails.
	Region string
	// KeepAlive to keep alive connection
	KeepAlive bool
	// MailerClient is the mailer client to use for sending emails.
	mailerClient MailerClient
}

// Mailer is a struct that holds the mailer instance and its configurations
type Mailer struct {
	host         string
	port         string
	username     string
	password     string
	apiService   APIServiceType
	apiKey       string
	emailToSend  chan Mail
	mailErr      chan error
	keepAlive    bool
	timeout      int
	mailerClient MailerClient
}

// NewMailer creates a new mailer instance.
func NewMailer(cfg MailCfg) *Mailer {
	mailer := &Mailer{
		host:         cfg.Host,
		port:         cfg.Port,
		username:     cfg.HostUser,
		password:     cfg.HostPassword,
		apiService:   cfg.APIService,
		apiKey:       cfg.APIKey,
		keepAlive:    cfg.KeepAlive,
		timeout:      cfg.Timeout,
		emailToSend:  make(chan Mail, 200),
		mailErr:      make(chan error),
		mailerClient: getMailerClient(cfg),
	}

	go mailer.listenForEmailsToBeSent()

	return mailer
}

// Send sends an email message using the chosen API service.
func (m *Mailer) Send(msg Mail) error {
	m.emailToSend <- msg
	return <-m.mailErr
}

// Close closes the emailToSend, result channels and the mailerClient.
func (m *Mailer) Close() {
	close(m.emailToSend)
	close(m.mailErr)
	m.mailerClient.Close()
}

// send sends the email message using the chosen API service.
func (m *Mailer) send(msg Mail) error {
	return m.mailerClient.Send(msg)
}

// ListenForEmailsToBeSent listens for email messages and sends them using the chosen API service.
// It is a blocking function that should be run in a goroutine.
func (m *Mailer) listenForEmailsToBeSent() {
	for msg := range m.emailToSend {
		err := m.send(msg)
		m.mailErr <- err
	}
}
