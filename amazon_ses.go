package mailer

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var (
	CharsetUtf8 = aws.String("UTF-8")
)

type sesParams struct {
	Region string
	Key    string
	Secret string
}

type sesMailer struct {
	sesClient *sesv2.Client
}

func newSES(params sesParams) MailerClient {
	creds := credentials.NewStaticCredentialsProvider(params.Key, params.Secret, "")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(params.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := sesv2.NewFromConfig(cfg)

	return &sesMailer{sesClient: client}
}

func (m *sesMailer) Send(msg Mail) error {
	mailInput := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data:   aws.String(msg.Html),
						Charset: CharsetUtf8,
					},
					Text: &types.Content{
						Data:   aws.String(msg.Text),
						Charset: CharsetUtf8,
					},
				},
				Subject: &types.Content{
					Data:    aws.String(msg.Subject),
					Charset: CharsetUtf8,
				},
			},
		},
		Destination:      &types.Destination{
			ToAddresses: []string{msg.To},
		},
		FromEmailAddress: aws.String(msg.From),
		ReplyToAddresses: []string{msg.ReplyTo},
	}

	_, err := m.sesClient.SendEmail(context.TODO(), mailInput)

	if err != nil {
		return err
	}

	return nil
}

func (m *sesMailer) Close() {

}
