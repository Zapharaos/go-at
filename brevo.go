package goat

import (
	"context"
	"encoding/base64"
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
func (s *BrevoService) Send(message *EmailMessage) error {
	brevoMsg := brevo.SendSmtpEmail{
		Sender:      s.from,
		To:          []brevo.SendSmtpEmailTo{{Email: message.To, Name: message.To}},
		Subject:     message.Subject,
		TextContent: message.PlainTextContent,
		HtmlContent: message.HTMLContent,
	}

	if message.ReplyTo != nil && message.ReplyTo.Address != "" {
		brevoMsg.ReplyTo = &brevo.SendSmtpEmailReplyTo{
			Name:  message.ReplyTo.Name,
			Email: message.ReplyTo.Address,
		}
	}

	if len(message.Headers) > 0 {
		h := make(map[string]interface{}, len(message.Headers))
		for k, v := range message.Headers {
			h[k] = v
		}
		brevoMsg.Headers = h
	}

	if len(message.Attachments) > 0 {
		atts := make([]brevo.SendSmtpEmailAttachment, 0, len(message.Attachments))
		for _, a := range message.Attachments {
			atts = append(atts, brevo.SendSmtpEmailAttachment{
				Content: base64.StdEncoding.EncodeToString(a.Content),
				Name:    a.Filename,
			})
		}
		brevoMsg.Attachment = atts
	}

	_, _, err := s.client.SendTransacEmail(context.Background(), brevoMsg)
	return err
}
