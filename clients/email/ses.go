package email

import (
	"context"
	"errors"

	"mbvlabs/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/smithy-go"
	"github.com/mbvlabs/andurel/pkg/postmand"
)

type SES struct {
	client *sesv2.Client
	set    string
}

func NewSES(cfg config.Mail) *SES {
	return &SES{
		client: sesv2.NewFromConfig(aws.Config{
			Region: cfg.SESRegion,
			Credentials: credentials.NewStaticCredentialsProvider(
				cfg.SESAccessKeyID, cfg.SESSecretAccessKey, "",
			),
		}),
		set: cfg.SESConfigurationSet,
	}
}

func (s *SES) SendTransactional(ctx context.Context, payload postmand.TransactionalPayload) error {
	input := emailInput(payload.From, payload.ReplyTo, payload.Subject, payload.HTMLBody, payload.TextBody)
	if payload.To != "" {
		input.Destination.ToAddresses = []string{payload.To}
	}
	input.Destination.CcAddresses = payload.Cc
	input.Destination.BccAddresses = payload.Bcc
	for _, attachment := range payload.Attachments {
		disposition := types.AttachmentContentDispositionAttachment
		if attachment.Inline {
			disposition = types.AttachmentContentDispositionInline
		}
		input.Content.Simple.Attachments = append(input.Content.Simple.Attachments, types.Attachment{
			FileName:                aws.String(attachment.Name),
			RawContent:              attachment.Content,
			ContentType:             aws.String(attachment.ContentType),
			ContentDisposition:      disposition,
			ContentTransferEncoding: types.AttachmentContentTransferEncodingBase64,
		})
	}
	return s.send(ctx, input, payload.Metadata)
}

func (s *SES) SendMarketing(ctx context.Context, payload postmand.MarketingPayload) error {
	if len(payload.To) == 0 {
		return postmand.ValidationError{Err: postmand.ErrMissingRecipient}
	}
	input := emailInput(payload.From, payload.ReplyTo, payload.Subject, payload.HTMLBody, payload.TextBody)
	// The existing job contract sends one marketing message per recipient.
	input.Destination.ToAddresses = payload.To[:1]
	for _, tag := range payload.Tags {
		input.EmailTags = append(input.EmailTags, types.MessageTag{Name: aws.String("tag"), Value: aws.String(tag)})
	}
	return s.send(ctx, input, payload.Metadata)
}

func emailInput(from, replyTo, subject, html, text string) *sesv2.SendEmailInput {
	message := &types.Message{Subject: &types.Content{Data: aws.String(subject)}, Body: &types.Body{}}
	if html != "" {
		message.Body.Html = &types.Content{Data: aws.String(html)}
	}
	if text != "" {
		message.Body.Text = &types.Content{Data: aws.String(text)}
	}
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(from),
		Destination:      &types.Destination{},
		Content:          &types.EmailContent{Simple: message},
	}
	if replyTo != "" {
		input.ReplyToAddresses = []string{replyTo}
	}
	return input
}

func (s *SES) send(ctx context.Context, input *sesv2.SendEmailInput, metadata map[string]string) error {
	if s.set != "" {
		input.ConfigurationSetName = aws.String(s.set)
	}
	for key, value := range metadata {
		input.EmailTags = append(input.EmailTags, types.MessageTag{Name: aws.String(key), Value: aws.String(value)})
	}
	_, err := s.client.SendEmail(ctx, input)
	if err == nil {
		return nil
	}
	if apiErr, ok := errors.AsType[smithy.APIError](err); ok {
		switch apiErr.ErrorCode() {
		case "BadRequestException", "ValidationError", "InvalidParameterValue", "MessageRejected":
			return postmand.ValidationError{Err: err}
		case "AccountSendingPausedException", "SendingPausedException", "MailFromDomainNotVerifiedException", "ConfigurationSetDoesNotExist", "NotFoundException":
			return postmand.PermanentError{Err: err}
		}
	}
	return postmand.TemporaryError{Err: err}
}
