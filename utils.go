package mailer

import (
	"fmt"
	"strings"

	mail "github.com/xhit/go-simple-mail/v2"
)

func getMailerClient(cfg MailCfg) (MailerClient, error) {
	if err := validateMailerRequiredFields(cfg); err != nil {
		return nil, err
	}

	switch cfg.APIService {
	case SMTP:
		return newSMTP(cfg.SMTP)
	case RESEND:
		return newResend(cfg.Resend), nil
	case AMAZON_SES:
		return newSES(cfg.SES)
	default:
		return nil, fmt.Errorf("invalid API service")
	}
}

func validateMailerRequiredFields(cfg MailCfg) error {
	switch cfg.APIService {
	case SMTP:
		if cfg.SMTP.Host == "" {
			return fmt.Errorf("missing required fields for SMTP i.e host")
		}
		if cfg.SMTP.Port == "" {
			return fmt.Errorf("missing required fields for SMTP i.e port")
		}
		if cfg.SMTP.Username == "" {
			return fmt.Errorf("missing required fields for SMTP i.e username")
		}
		if cfg.SMTP.Password == "" {
			return fmt.Errorf("missing required fields for SMTP i.e password")
		}
	case RESEND:
		if cfg.Resend.apiKey == "" {
			return fmt.Errorf("API key is missing")
		}
	case AMAZON_SES:
		if cfg.SES.Key == "" {
			return fmt.Errorf("API key is missing")
		}
		if cfg.SES.Secret == "" {
			return fmt.Errorf("API secret is missing")
		}
		if cfg.SES.Region == "" {
			return fmt.Errorf("region is missing")
		}
	}
	return nil
}

func getSplitEmails(emails string) []string {
	if emails == "" {
		return []string{}
	}
	return strings.Split(emails, ",")
}

func buildMessage(msg Mail) (string, *mail.Email) {
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	if msg.Text != "" {
		email.SetBody(mail.TextPlain, msg.Text)
	}

	if msg.Html != "" {
		email.AddAlternative(mail.TextHTML, msg.Html)
	}

	if msg.ReplyTo != "" {
		email.SetReplyTo(msg.ReplyTo)
	}
	if msg.Cc != "" {
		email.AddCc(msg.Cc)
	}
	if msg.Bcc != "" {
		email.AddBcc(msg.Bcc)
	}

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.Attach(&mail.File{
				FilePath: attachment.Path,
				Name:     attachment.Name,
				Inline:   true,
			})
		}
	}

	return email.GetMessage(), email
}
