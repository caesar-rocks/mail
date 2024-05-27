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

func newResend(params resendParams) Mailer {
	client := resend.NewClient(params.apiKey)

	return &resendMailer{
		resendClient: client,
	}
}

func (m *resendMailer) Send(msg MailerMessage) error {
	params := &resend.SendEmailRequest{
		To:      getSplitEmails(msg.To),
		From:    msg.From,
		Text:    msg.Text,
		Subject: msg.Subject,
		Cc:      getSplitEmails(msg.Cc),
		Bcc:     getSplitEmails(msg.Bcc),
		ReplyTo: msg.ReplyTo,
	}

	_, err := m.resendClient.Emails.Send(params)
	if err != nil {
		panic(err)
	}
	return nil
}

func (m *resendMailer) Close() {
	// Not implemented because resend-go does not have a Close method.
}
