package awsmail

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/danikarik/okpock/pkg/mail"
)

// New returns AWS SES Mailer.
func New(region string) (mail.Mailer, error) {
	cfg := &aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Region:      aws.String(region),
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create session: %v", err)
	}
	mailer := &sesMailer{
		srv: ses.New(sess),
	}
	return mailer, nil
}

type sesMailer struct {
	srv *ses.SES
}

func (s *sesMailer) createInput(message *mail.Message) *ses.SendEmailInput {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(message.Recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(message.Charset),
					Data:    aws.String(message.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(message.Charset),
					Data:    aws.String(message.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(message.Charset),
				Data:    aws.String(message.Subject),
			},
		},
		Source: aws.String(message.Sender),
	}
	return input
}

func (s *sesMailer) SendMail(ctx context.Context, message *mail.Message) (*time.Time, error) {
	err := message.Valid()
	if err != nil {
		return nil, err
	}

	input := s.createInput(message)
	output, err := s.srv.SendEmail(input)
	if err != nil {
		return nil, fmt.Errorf("could not send email: %v", err)
	}

	message.ID = *output.MessageId
	now := time.Now()

	return &now, nil
}
