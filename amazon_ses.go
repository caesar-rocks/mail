package mailer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var (
	CharsetUtf8 = aws.String("UTF-8")
)

type sesMailerClient interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

type SESCfg struct {
	Region string
	Key    string
	Secret string
}

type sesMailer struct {
	sesClient sesMailerClient
}

func newSES(params SESCfg) (MailerClient, error) {
	creds := credentials.NewStaticCredentialsProvider(params.Key, params.Secret, "")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(params.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	client := sesv2.NewFromConfig(cfg)

	return &sesMailer{sesClient: client}, nil
}

func (m *sesMailer) Send(msg Mail) error {
	message, _ := buildMessage(msg)
	mailInput := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: []byte(message),
			},
		},
	}

	_, err := m.sesClient.SendEmail(context.TODO(), mailInput)

	if err != nil {
		return err
	}

	return nil
}

func (m *sesMailer) Close() {
	// No need to close the connection
}
