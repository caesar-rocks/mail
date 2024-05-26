package mailer

type Attachment struct {
	// Name is the name of the attachment.
	Name string
	// Path is the path to the attachment.
	Path string
}

type MailerMessage struct {
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

type Mailer interface {
	Send(msg MailerMessage) error
	Close()
}
