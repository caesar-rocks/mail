package mailer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	mail "github.com/xhit/go-simple-mail/v2"
)

func getMailerClient(cfg MailCfg) MailerClient {
	err := validateMailerRequiredFields(cfg)
	if err != nil {
		panic(err)
	}
	if cfg.mailerClient != nil {
		return cfg.mailerClient
	}
	switch cfg.APIService {
	case SMTP:
		return newSMTP(smtpParams{
			Host:      cfg.Host,
			Port:      cfg.Port,
			Username:  cfg.HostUser,
			Password:  cfg.HostPassword,
			KeepAlive: cfg.KeepAlive,
			Timeout:   cfg.Timeout,
			useTLS:    cfg.UseTLS,
		})
	case RESEND:
		return newResend(resendParams{
			apiKey: cfg.APIKey,
		})
	case SENDGRID, MAILGUN:
		return nil
	case AMAZON_SES:
		return newSES(
			sesParams{
				Region: cfg.Region,
				Key:    cfg.APIKey,
				Secret: cfg.APISecret,
			},
		)
	default:
		panic("invalid API service")
	}
}

func getPort(port string) int {
	p, err := strconv.Atoi(port)

	if err != nil {
		log.Fatalln(err)
	}

	return p
}

func validateMailerRequiredFields(cfg MailCfg) error {
	switch cfg.APIService {
	case SMTP:
		if cfg.Host == "" || cfg.Port == "" || cfg.HostUser == "" || cfg.HostPassword == "" {
			return fmt.Errorf("missing required fields for SMTP i.e host, port, username, password")
		}
	case SENDGRID, MAILGUN, RESEND:
		if cfg.APIKey == "" {
			return fmt.Errorf("API key is missing")
		}
	case AMAZON_SES:
		if cfg.APIKey == "" || cfg.APISecret == "" || cfg.Region == "" {
			return fmt.Errorf("missing required fields for Amazon SES i.e region, key, secret")
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
