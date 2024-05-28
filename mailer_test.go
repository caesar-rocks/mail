package mailer

import (
	"testing"
)

const (
	MailHost         = "testhost"
	MailPort         = "1234"
	MailHostUser     = "testuser"
	MailHostPassword = "testpassword"
	MailAPIService   = "smtp"
	MailAPIKey       = "testapikey"
	MailAPISecret	= "testapisecret"
	MailRegion		= "us-west-2"
)

func TestMail_NewMail(t *testing.T) {
	mockClient := &mockMailerClient{}

	testCases := []struct {
		name    string
		success bool
		cfg     MailCfg
	}{
		{
			name:    "Should create new instance of new mailer with smtp fields",
			success: true,
			cfg: MailCfg{
				Host:         MailHost,
				Port:         MailPort,
				HostUser:     MailHostUser,
				HostPassword: MailHostPassword,
				APIService:   MailAPIService,
				APIKey:       MailAPIKey,
				mailerClient: mockClient,
			},
		},
		{
			name:    "Should not create new instance of new mailer with missing fields",
			success: false,
			cfg: MailCfg{
				Host:         MailHost,
				Port:         MailPort,
				HostUser:     MailHostUser,
				HostPassword: "",
				APIService:   "",
				APIKey:       MailAPIKey,
				mailerClient: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.success {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected NewMail to panic with incomplete configuration")
					}
				}()
			}
			mailer := NewMailer(tc.cfg)
			if tc.success {
				if mailer == nil {
					t.Errorf("Expected mailer to be created, got nil")
				}
				mailer.Close()
			}
		})

	}
}

func TestNewMail_SendMail(t *testing.T) {
	mockClient := &mockMailerClient{}

	testCases := []struct {
		name       string
		apiService APIServiceType
	}{
		{
			name:       "send email with smtp",
			apiService: SMTP,
		},
		{
			name:       "send email with resend",
			apiService: RESEND,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mailer := NewMailer(MailCfg{
				Host:         MailHost,
				Port:         MailPort,
				HostUser:     MailHostUser,
				HostPassword: MailHostPassword,
				APIService:   tc.apiService,
				APIKey:       MailAPIKey,
				mailerClient: mockClient,
			})

			err := mailer.Send(Mail{
				To:      "test@example.com",
				Subject: "Test",
				Html:    "<h1>Test</h1>",
				Text:    "Test",
			})

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

		})
	}
}

type mockMailerClient struct {
}

func (m *mockMailerClient) Send(msg Mail) error {
	return nil
}

func (m *mockMailerClient) Close() {

}
