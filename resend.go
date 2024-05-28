package mailer

import (
	"github.com/resend/resend-go/v2"
)

type resendParams struct {
	apiKey string
}

type resendMailer struct {
	resendClient *resend.Client
}

func newResend(params resendParams) MailerClient {
	client := resend.NewClient(params.apiKey)

	return &resendMailer{resendClient: client}
}

func (m *resendMailer) Send(msg Mail) error {
	params := &resend.SendEmailRequest{
		To:          getSplitEmails(msg.To),
		From:        msg.From,
		Text:        msg.Text,
		Html:        msg.Html,
		Attachments: m.getAttachments(msg.Attachments),
		Subject:     msg.Subject,
		Cc:          getSplitEmails(msg.Cc),
		Bcc:         getSplitEmails(msg.Bcc),
		ReplyTo:     msg.ReplyTo,
	}

	_, err := m.resendClient.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}

func (m *resendMailer) getAttachments(a []Attachment) []*resend.Attachment {
	var attachments []*resend.Attachment
	for _, attachment := range a {
		attachments = append(attachments, &resend.Attachment{
			Filename: attachment.Name,
			Path:     attachment.Path,
		})
	}
	return attachments
}

func (m *resendMailer) Close() {
	// Not implemented because resend-go does not have a Close method.
}
