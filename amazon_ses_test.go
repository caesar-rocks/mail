package mailer

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func TestAWSSesClient(t *testing.T) {

	testCases := []struct {
		name string
		cfg  sesParams
	}{
		{
			name: "Should create new instance of new mailer with ses fields",
			cfg: sesParams{
				Region: "us-west-2",
				Key:    "testkey",
				Secret: "testsecret",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ses := newSES(tc.cfg)
			if ses == nil {
				t.Errorf("Expected newSES to return a mailer, got nil")
			} else {
				t.Logf("SES mailer created successfully")
			}
		})
	}
}

func TestAWSSes_Send(t *testing.T) {
	testCases := []struct {
		name    string
		success bool
		payload Mail
	}{
		{
			name: "Should send email successfully",
			payload: Mail{
				Subject: "test",
				From:    "info@test.com",
				To:      "test@gmail.com",
				Html:    "<p>test</p>",
				Text:    "test",
			},
			success: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ses := newSES(sesParams{
				Region: "us-west-2",
				Key:    "testkey",
				Secret: "testsecret",
			})

			mockClient := &mockSESClient{
				SendEmailFunc: func(ctx context.Context, input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error) {
					if tc.success {
						return &sesv2.SendEmailOutput{
							MessageId: aws.String("test-id"),
						}, nil
					} else {
						return nil, errors.New("error")
					}

				},
			}

			ses.(*sesMailer).sesClient = mockClient

			err := ses.Send(tc.payload)

			if tc.success {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			}

		})
	}
}

type mockSESClient struct {
	SendEmailFunc func(ctx context.Context, input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error)
}

func (m *mockSESClient) SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	return m.SendEmailFunc(ctx, params)
}
