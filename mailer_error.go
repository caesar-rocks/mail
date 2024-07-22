package mailer

import "fmt"

type MailerError struct {
	API        string
	StatusCode int
	Meta       map[string]interface{}
}

func (e *MailerError) Error() string {
	return fmt.Sprintf(`{"API": "%s", "StatusCode": %d, "Meta": %v}`, e.API, e.StatusCode, e.Meta)
}

func newMailerError(api string, statusCode int, meta map[string]interface{}) error {
	return &MailerError{
		API:        api,
		StatusCode: statusCode,
		Meta:       meta,
	}
}
