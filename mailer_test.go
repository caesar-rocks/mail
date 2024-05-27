package mailer

import (
	"os"
	"sync"
	"testing"
)

func TestMail_NewMail(t *testing.T) {
	mockClient := &mockMailerClient{}

	os.Setenv("MAIL_HOST", "testhost")
	os.Setenv("MAIL_PORT", "1234")
	os.Setenv("MAIL_HOST_USER", "testuser")
	os.Setenv("MAIL_HOST_PASSWORD", "testpassword")
	os.Setenv("MAIL_API_SERVICE", "smtp")
	os.Setenv("MAIL_API_KEY", "testapikey")

	testCases := []struct {
		name    string
		success bool
		cfg     MailConfig
	}{
		{
			name:    "Should create new instance of new mailer with all fields",
			success: true,
			cfg: MailConfig{
				Host:         os.Getenv("MAIL_HOST"),
				Port:         os.Getenv("MAIL_PORT"),
				HostUser:     os.Getenv("MAIL_HOST_USER"),
				HostPassword: os.Getenv("MAIL_HOST_PASSWORD"),
				APIService:   "smtp",
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: mockClient,
			},
		},
		{
			name:    "Should not create new instance of new mailer with missing fields",
			success: false,
			cfg: MailConfig{
				Host:         os.Getenv("MAIL_HOST"),
				Port:         os.Getenv("MAIL_PORT"),
				HostUser:     os.Getenv("MAIL_HOST_USER"),
				HostPassword: "",
				APIService:   "",
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetMailer()
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

	os.Setenv("MAIL_HOST", "smtp.example.com")
	os.Setenv("MAIL_PORT", "587")
	os.Setenv("MAIL_HOST_USER", "info@caesar.rock")
	os.Setenv("MAIL_HOST_PASSWORD", "supersecretpassword")

	testCases := []struct {
		name       string
		apiService APIServiceType
	}{
		{
			name:       "send email with smtp",
			apiService: SMTP,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetMailer()

			mailer := NewMailer(MailConfig{
				Host:         os.Getenv("MAIL_HOST"),
				Port:         os.Getenv("MAIL_PORT"),
				HostUser:     os.Getenv("MAIL_HOST_USER"),
				HostPassword: os.Getenv("MAIL_HOST_PASSWORD"),
				APIService:   tc.apiService,
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: mockClient,
			})

			err := mailer.Send(MailerMessage{
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

func (m *mockMailerClient) Send(msg MailerMessage) error {
	return nil
}

func (m *mockMailerClient) Close() {

}

func resetMailer() {
	mailer = nil
	once = sync.Once{}
}
