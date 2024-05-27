package mailer

import (
	"testing"
)

func TestGetMailer(t *testing.T) {
	mockClient := &mockMailerClient{}

	testCases := []struct {
		name       string
		apiService APIServiceType
		cfg        MailCfg
		success    bool
	}{
		{
			name:       "get smtp",
			apiService: SMTP,
			cfg: MailCfg{
				Host:         MailHost,
				Port:         MailPort,
				HostUser:     MailHostUser,
				HostPassword: MailHostPassword,
				APIService:   SMTP,
				APIKey:       MailAPIKey,
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get sendgrid",
			apiService: SENDGRID,
			cfg: MailCfg{
				APIService:   SENDGRID,
				APIKey:       MailAPIKey,
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get mailgun",
			apiService: "mailgun",
			cfg: MailCfg{
				APIService:   MAILGUN,
				APIKey:       MailAPIKey,
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get amazon_ses",
			apiService: AMAZON_SES,
			cfg: MailCfg{
				APIService:   AMAZON_SES,
				APIKey:       MailAPIKey,
				mailerClient: mockClient,
			},
			success: true,
		},
		{
			name:       "get unknown",
			apiService: "unknown",
			cfg: MailCfg{
				APIService:   "unknown",
				APIKey:       MailAPIKey,
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
		cfg     MailCfg
		success bool
	}{
		{
			name: "validate smtp with missing username",
			cfg: MailCfg{
				Host:         MailHost,
				Port:         MailPort,
				HostUser:     "",
				HostPassword: MailHostPassword,
				APIService:   SMTP,
				APIKey:       MailAPIKey,
			},
			success: false,
		},
		{
			name: "validate smtp",
			cfg: MailCfg{
				Host:         MailHost,
				Port:         MailPort,
				HostUser:     MailHostUser,
				HostPassword: MailHostPassword,
				APIService:   SMTP,
				APIKey:       MailAPIKey,
			},
			success: true,
		},
		{
			name: "validate resend",
			cfg: MailCfg{
				APIService: RESEND,
				APIKey:     "",
			},
			success: false,
		},
		{
			name: "validate sendgrid",
			cfg: MailCfg{
				APIService: SENDGRID,
				APIKey:     MailAPIKey,
			},
			success: true,
		},
		{
			name: "validate mailgun",
			cfg: MailCfg{
				APIService: MAILGUN,
				APIKey:     MailAPIKey,
			},
			success: true,
		},
		{
			name: "validate amazon_ses",
			cfg: MailCfg{
				APIService: AMAZON_SES,
				APIKey:     MailAPIKey,
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
