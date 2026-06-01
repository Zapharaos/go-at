package goat

import (
	"context"
	"net/http"

	brevo "github.com/getbrevo/brevo-go/lib"
)

// BrevoClient is an interface for sending emails
type BrevoClient interface {
	SendTransacEmail(ctx context.Context, sendSmtpEmail brevo.SendSmtpEmail) (brevo.CreateSmtpEmail, *http.Response, error)
}

// BrevoService implements the SenderService interface using Brevo
type BrevoService struct {
	client BrevoClient
	from   *brevo.SendSmtpEmailSender
}

// NewBrevoService returns a new instance of BrevoService
func NewBrevoService(apiKey, senderName, senderEmail string) SenderService {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)

	s := BrevoService{
		client: brevo.NewAPIClient(cfg).TransactionalEmailsApi,
		from: &brevo.SendSmtpEmailSender{
			Name:  senderName,
			Email: senderEmail,
		},
	}
	var service SenderService = &s
	return service
}

// Send sends an email using Brevo
func (s *BrevoService) Send(emailTo, subject, plainTextContent, htmlContent string) error {
	message := brevo.SendSmtpEmail{
		Sender:      s.from,
		To:          []brevo.SendSmtpEmailTo{{Email: emailTo, Name: emailTo}},
		Subject:     subject,
		TextContent: plainTextContent,
		HtmlContent: htmlContent,
	}

	// Send email
	_, _, err := s.client.SendTransacEmail(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}

// SendWithReplyTo sends an email using Brevo with a reply-to address
func (s *BrevoService) SendWithReplyTo(emailTo, subject, plainTextContent, htmlContent string, replyTo ReplyTo) error {
	message := brevo.SendSmtpEmail{
		Sender:      s.from,
		To:          []brevo.SendSmtpEmailTo{{Email: emailTo, Name: emailTo}},
		Subject:     subject,
		TextContent: plainTextContent,
		HtmlContent: htmlContent,
	}
	if replyTo.Address != "" {
		message.ReplyTo = &brevo.SendSmtpEmailReplyTo{
			Name:  replyTo.Name,
			Email: replyTo.Address,
		}
	}

	// Send email
	_, _, err := s.client.SendTransacEmail(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
