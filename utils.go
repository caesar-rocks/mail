package mailer

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func getMailerClient(cfg MailConfig) Mailer {
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
	case SENDGRID, MAILGUN:
		return nil
	case AMAZON_SES:
		return nil
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

func validateMailerRequiredFields(cfg MailConfig) error {
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
		return nil
	}
	return nil
}

func getSplitEmails(emails string) []string {
	if emails == "" {
		return []string{}
	}
	return strings.Split(emails, ",")
}
