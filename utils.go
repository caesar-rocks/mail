package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	case POSTMARK:
		return newPostmark(cfg.Postmark), nil
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
		if cfg.Resend.APIKey == "" {
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

func makeApiCall(endpoint string, headers map[string]string, body map[string]any, httpClient *http.Client) error {
	prepareBody, err := prepareBody(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, prepareBody)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	err = checkResponse(res)

	if err != nil {
		return err
	}

	return nil
}

func prepareBody(body map[string]any) (io.Reader, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBytes), nil
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 400 {
		metaDecoder := json.NewDecoder(res.Body)
		var meta map[string]interface{}
		err := metaDecoder.Decode(&meta)
		if err != nil {
			return err
		}
		return newMailerError("Mailer API", res.StatusCode, meta)
	}
	return nil
}
