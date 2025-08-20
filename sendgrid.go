package goat

import (
	"context"

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
func (s *SendgridService) Send(emailTo, subject, plainTextContent, htmlContent string) error {
	// Email props
	to := mail.NewEmail(emailTo, emailTo)
	message := mail.NewSingleEmail(s.from, subject, to, plainTextContent, htmlContent)

	// Send email
	_, err := s.client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
