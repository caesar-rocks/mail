package mailer

import (
	"encoding/base64"
	"fmt"
	"mime"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendgridMailer struct {
	sendGridClient *sendgrid.Client
}

type sendGridParams struct {
	apiKey string
}

func newSendgrid(params sendGridParams) MailerClient {
	client := sendgrid.NewSendClient(params.apiKey)

	return &sendgridMailer{sendGridClient: client}
}

func (m *sendgridMailer) Send(msg Mail) error {
	from := mail.NewEmail(msg.From, msg.From)
	to := mail.NewEmail(msg.To, msg.To)

	message := mail.NewSingleEmail(from, msg.Subject, to, msg.Text, msg.Html)
	if msg.Cc != "" {
		message.Personalizations[0].AddCCs(mail.NewEmail(msg.Cc, msg.Cc))
	}
	if msg.Bcc != "" {
		message.Personalizations[0].AddBCCs(mail.NewEmail(msg.Bcc, msg.Bcc))
	}

	if msg.ReplyTo != "" {
		message.ReplyTo = mail.NewEmail(msg.ReplyTo, msg.ReplyTo)
	}

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			file, err := os.ReadFile(attachment.Path)
			if err != nil {
				return fmt.Errorf("error reading attachment file: %v", err)
			}
			encoded := base64.StdEncoding.EncodeToString([]byte(file))

			message.AddAttachment(&mail.Attachment{
				Content:     encoded,
				Filename:    attachment.Name,
				Type:        mime.TypeByExtension(attachment.Path),
				Name:        attachment.Name,
				Disposition: "attachment",
			})
		}
	}

	_, err := m.sendGridClient.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func (m *sendgridMailer) Close() {
	// do nothing
}
