package mailer

import (
	"encoding/base64"
	"io"
	"mime"
	"net/http"
	"os"
	"time"
)

const (
	postmarkAPIEndpoint = "https://api.postmarkapp.com/email"
)

type PostmarkCfg struct {
	ServerToken string
}

type postmarkMailer struct {
	serverToken string
	client      *http.Client
}

type postAttachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
}

func newPostmark(params PostmarkCfg) MailerClient {
	return &postmarkMailer{
		serverToken: params.ServerToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (m *postmarkMailer) Send(msg Mail) error {
	headers := map[string]string{
		"Accept":                  "application/json",
		"Content-Type":            "application/json",
		"X-Postmark-Server-Token": m.serverToken,
	}
	var attachments []postAttachment
	for _, attachment := range msg.Attachments {
		base64Content, err := getBase64Content(attachment.Path)
		if err != nil {
			return newMailerError("Postmark", 500, map[string]interface{}{
				"error": "Error getting base64 content",
			})
		}
		attachments = append(attachments, postAttachment{
			Name:        attachment.Name,
			Content:     base64Content,
			ContentType: mime.TypeByExtension(attachment.Name),
		})

	}
	body := map[string]any{
		"From":          msg.From,
		"To":            msg.To,
		"Subject":       msg.Subject,
		"HtmlBody":      msg.Html,
		"TextBody":      msg.Text,
		"Cc":            msg.Cc,
		"Bcc":           msg.Bcc,
		"ReplyTo":       msg.ReplyTo,
		"Attachments":   attachments,
		"Headers":       msg.Headers,
		"MessageStream": "outbound",
	}
	return makeApiCall(postmarkAPIEndpoint, headers, body, m.client)
}

func (m *postmarkMailer) Close() {
	m.client.CloseIdleConnections()
}

func getBase64Content(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), err
}
