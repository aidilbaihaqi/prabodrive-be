package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/aidilbaihaqi/prabodrive-be/internal/config"
)

type EmailService struct {
	client    *ses.Client
	fromEmail string
}

func NewEmailService(sesClient *ses.Client, cfg config.SESConfig) *EmailService {
	return &EmailService{
		client:    sesClient,
		fromEmail: cfg.FromEmail,
	}
}

func (e *EmailService) SendShareNotification(ctx context.Context, toEmail, docName, shareURL string) error {
	subject := fmt.Sprintf("Document shared: %s", docName)
	body := fmt.Sprintf(
		"A share link has been created for document \"%s\".\n\nAccess it here: %s",
		docName, shareURL,
	)

	_, err := e.client.SendEmail(ctx, &ses.SendEmailInput{
		Source: aws.String(e.fromEmail),
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Message: &types.Message{
			Subject: &types.Content{Data: aws.String(subject)},
			Body: &types.Body{
				Text: &types.Content{Data: aws.String(body)},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses: send email: %w", err)
	}
	return nil
}
