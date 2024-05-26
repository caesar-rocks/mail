package mailer

import (
	"fmt"
	"log"
	"strconv"
)

func getMailerClient(opt MailConfig) Mailer {
	err := validateMailerRequiredFields(opt)
	if err != nil {
		panic(err)
	}
	if opt.mailerClient != nil {
		return opt.mailerClient
	}
	switch opt.APIService {
	case SMTP:
		return newSMTP(smtpParams{
			Host:      opt.Host,
			Port:      opt.Port,
			Username:  opt.HostUser,
			Password:  opt.HostPassword,
			KeepAlive: opt.KeepAlive,
			Timeout:   opt.Timeout,
			useTLS:    opt.UseTLS,
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

func validateMailerRequiredFields(opt MailConfig) error {
	switch opt.APIService {
	case SMTP:
		if opt.Host == "" || opt.Port == "" || opt.HostUser == "" || opt.HostPassword == "" {
			return fmt.Errorf("SMTP configuration is missing")
		}
	case SENDGRID, MAILGUN, RESEND:
		if opt.APIKey == "" {
			return fmt.Errorf("API key is missing")
		}
	case AMAZON_SES:
		return nil
	}
	return nil
}
