package goat

import (
	"context"
	"os"

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
func NewSendgridService() SenderService {
	// Get SendGrid API key and sender info from environment variables
	apiKey := os.Getenv("SENDGRID_API_KEY")
	senderName := os.Getenv("SENDGRID_SENDER_NAME")
	senderEmail := os.Getenv("SENDGRID_SENDER_EMAIL")

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
