package mailer

type APIServiceType string

const (
	SMTP       APIServiceType = "smtp"
	AMAZON_SES APIServiceType = "amazon-ses"
	RESEND     APIServiceType = "resend"
	POSTMARK   APIServiceType = "postmark"
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
	// Headers is a string of headers.
	Headers map[string]string
	// Attachments is an array of attachments.
	Attachments []Attachment
}

type MailerClient interface {
	Send(msg Mail) error
	Close()
}

// MailCfg is a struct that holds the configuration for the mailer.
type MailCfg struct {
	// The configuration for the SMTP server.
	SMTP SMTPCfg

	// The configuration for the Amazon SES service.
	SES SESCfg

	// The configuration for the Resend service.
	Resend ResendCfg

	// The configuration for the Postmark service
	Postmark PostmarkCfg

	// APIService is the service to use for sending emails.
	APIService APIServiceType
}

// Mailer is a struct that holds the mailer instance and its configurations
type Mailer struct {
	emailToSend  chan Mail
	mailErr      chan error
	mailerClient MailerClient
}

// NewMailer creates a new mailer instance.
func NewMailer(cfg MailCfg) (*Mailer, error) {
	mailerClient, err := getMailerClient(cfg)
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		emailToSend:  make(chan Mail, 200),
		mailErr:      make(chan error),
		mailerClient: mailerClient,
	}

	go mailer.listenForEmailsToBeSent()

	return mailer, nil
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
