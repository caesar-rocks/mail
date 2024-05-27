package mailer

import (
	"os"
	"testing"
)

func TestGetMailer(t *testing.T) {
	mockClient := &mockMailerClient{}

	testCases := []struct {
		name       string
		apiService APIServiceType
		cfg        MailConfig
		success    bool
	}{
		{
			name:       "get smtp",
			apiService: SMTP,
			cfg: MailConfig{
				Host:         os.Getenv("MAIL_HOST"),
				Port:         os.Getenv("MAIL_PORT"),
				HostUser:     os.Getenv("MAIL_HOST_USER"),
				HostPassword: os.Getenv("MAIL_HOST_PASSWORD"),
				APIService:   SMTP,
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get sendgrid",
			apiService: SENDGRID,
			cfg: MailConfig{
				APIService:   SENDGRID,
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get mailgun",
			apiService: "mailgun",
			cfg: MailConfig{
				APIService:   MAILGUN,
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get amazon_ses",
			apiService: AMAZON_SES,
			cfg: MailConfig{
				APIService:   AMAZON_SES,
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get unknown",
			apiService: "unknown",
			cfg: MailConfig{
				APIService:   "unknown",
				APIKey:       os.Getenv("MAIL_API_KEY"),
				mailerClient: nil,
			},
			success: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.success {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected getMailerClient to panic with unknown api service")
					}
				}()
			}
			mailerClient := getMailerClient(tc.cfg)
			if tc.success {
				if mailerClient == nil {
					t.Errorf("Expected getMailerClient return a mailer, got nil")
				}
			}
		})
	}
}

func TestGetPort(t *testing.T) {
	testCases := []struct {
		name     string
		port     string
		expected int
	}{
		{
			name:     "get port",
			port:     "587",
			expected: 587,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			port := getPort(tc.port)
			if port != tc.expected {
				t.Errorf("Expected port to be %d, got %d", tc.expected, port)
			}
		})
	}
}

func TestValidateMailerRequiredFields(t *testing.T) {
	testCases := []struct {
		name    string
		cfg     MailConfig
		success bool
	}{
		{
			name: "validate smtp",
			cfg: MailConfig{
				Host:         os.Getenv("MAIL_HOST"),
				Port:         os.Getenv("MAIL_PORT"),
				HostUser:     "",
				HostPassword: os.Getenv("MAIL_HOST_PASSWORD"),
				APIService:   SMTP,
				APIKey:       os.Getenv("MAIL_API_KEY"),
			},
			success: false,
		},
		{
			name: "validate smtp",
			cfg: MailConfig{
				Host:         os.Getenv("MAIL_HOST"),
				Port:         os.Getenv("MAIL_PORT"),
				HostUser:     os.Getenv("MAIL_HOST_USER"),
				HostPassword: os.Getenv("MAIL_HOST_PASSWORD"),
				APIService:   SMTP,
				APIKey:       os.Getenv("MAIL_API_KEY"),
			},
			success: true,
		},
		{
			name: "validate resend",
			cfg: MailConfig{
				APIService: RESEND,
				APIKey:     "",
			},
			success: false,
		},
		{
			name: "validate sendgrid",
			cfg: MailConfig{
				APIService: SENDGRID,
				APIKey:     os.Getenv("MAIL_API_KEY"),
			},
			success: true,
		},
		{
			name: "validate mailgun",
			cfg: MailConfig{
				APIService: MAILGUN,
				APIKey:     os.Getenv("MAIL_API_KEY"),
			},
			success: true,
		},
		{
			name: "validate amazon_ses",
			cfg: MailConfig{
				APIService: AMAZON_SES,
				APIKey:     os.Getenv("MAIL_API_KEY"),
			},
			success: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMailerRequiredFields(tc.cfg)
			if tc.success {
				if err != nil {
					t.Errorf("Expected validateMailerRequiredFields to return nil, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected validateMailerRequiredFields to return error, got nil")
				}
			}
		})
	}
}

func TestGetSplitEmails(t *testing.T) {
	testCases := []struct {
		name     string
		emails   string
		expected []string
	}{
		{
			name:     "get split emails - single",
			emails:   "test1@gmail.com,test2@gmail.com",
			expected: []string{"test1@gmail.com", "test2@gmail.com"},
		},
		{
			name:     "get split emails - empty",
			emails:   "",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emails := getSplitEmails(tc.emails)
			if len(emails) != len(tc.expected) {
				t.Errorf("Expected emails to be %v, got %v", tc.expected, emails)
			}
		})
	}
}
