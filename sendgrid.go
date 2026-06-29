package goat

import (
	"context"
	"encoding/base64"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendgridClient is an interface for sending emails
type SendgridClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
	SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error)
}

// SendgridService implements the SenderService interface using SendGrid
type SendgridService struct {
	client SendgridClient
	from   *mail.Email
}

// NewSendgridService returns a new instance of SendgridService
func NewSendgridService(apiKey, senderName, senderEmail string) SenderService {
	s := SendgridService{
		client: sendgrid.NewSendClient(apiKey),
		from:   mail.NewEmail(senderName, senderEmail),
	}
	var service SenderService = &s
	return service
}

// Send sends an email using SendGrid
func (s *SendgridService) Send(message *EmailMessage) error {
	to := mail.NewEmail(message.To, message.To)
	msg := mail.NewSingleEmail(s.from, message.Subject, to, message.PlainTextContent, message.HTMLContent)

	if message.ReplyTo != nil && message.ReplyTo.Address != "" {
		msg.SetReplyTo(mail.NewEmail(message.ReplyTo.Name, message.ReplyTo.Address))
	}

	for k, v := range message.Headers {
		msg.SetHeader(k, v)
	}

	for _, a := range message.Attachments {
		att := mail.NewAttachment()
		att.SetContent(base64.StdEncoding.EncodeToString(a.Content))
		att.SetType(a.ContentType)
		att.SetFilename(a.Filename)
		if a.ContentID != "" {
			att.SetContentID(a.ContentID)
			att.SetDisposition("inline")
		} else {
			att.SetDisposition("attachment")
		}
		msg.AddAttachment(att)
	}

	_, err := s.client.Send(msg)
	return err
}
