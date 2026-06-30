package goat

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/stretchr/testify/assert"
)

type MockBrevoClient struct {
	SendResponse brevo.CreateSmtpEmail
	SendHTTP     *http.Response
	SendError    error
	LastEmail    brevo.SendSmtpEmail
}

func (m *MockBrevoClient) SendTransacEmail(ctx context.Context, sendSmtpEmail brevo.SendSmtpEmail) (brevo.CreateSmtpEmail, *http.Response, error) {
	m.LastEmail = sendSmtpEmail
	return m.SendResponse, m.SendHTTP, m.SendError
}

// TestNewBrevoService tests the NewBrevoService function
func TestNewBrevoService(t *testing.T) {
	service := NewBrevoService("test_api_key", "test_sender_name", "test_sender_email")

	assert.NotNil(t, service)
	assert.IsType(t, &BrevoService{}, service)
}

// TestBrevoService_Send tests the Send method of BrevoService
func TestBrevoService_Send(t *testing.T) {
	service := NewBrevoService("test_api_key", "test_sender_name", "test_sender_email")

	t.Run("Success - plain message", func(t *testing.T) {
		service.(*BrevoService).client = &MockBrevoClient{SendError: nil}

		err := service.Send(NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content"))
		assert.NoError(t, err)
	})

	t.Run("Success - with reply-to", func(t *testing.T) {
		service.(*BrevoService).client = &MockBrevoClient{SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with headers", func(t *testing.T) {
		service.(*BrevoService).client = &MockBrevoClient{SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>, <https://example.com/unsubscribe>")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with reply-to and headers combined", func(t *testing.T) {
		service.(*BrevoService).client = &MockBrevoClient{SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with attachment", func(t *testing.T) {
		mock := &MockBrevoClient{SendError: nil}
		service.(*BrevoService).client = mock

		content := []byte("PDF bytes")
		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithAttachment("invoice.pdf", "application/pdf", content)
		err := service.Send(msg)
		assert.NoError(t, err)

		assert.Len(t, mock.LastEmail.Attachment, 1)
		att := mock.LastEmail.Attachment[0]
		assert.Equal(t, "invoice.pdf", att.Name)
		assert.Equal(t, base64.StdEncoding.EncodeToString(content), att.Content)
	})

	t.Run("Success - inline attachment degrades to regular attachment", func(t *testing.T) {
		mock := &MockBrevoClient{SendError: nil}
		service.(*BrevoService).client = mock

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "<img src=\"cid:logo\">").
			WithInlineAttachment("logo.png", "image/png", []byte("PNG bytes"), "logo")
		err := service.Send(msg)
		assert.NoError(t, err)

		// Brevo SDK has no Content-ID field; the attachment is still sent by name.
		assert.Len(t, mock.LastEmail.Attachment, 1)
		assert.Equal(t, "logo.png", mock.LastEmail.Attachment[0].Name)
	})

	t.Run("Success - attachment with reply-to and headers combined", func(t *testing.T) {
		mock := &MockBrevoClient{SendError: nil}
		service.(*BrevoService).client = mock

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>").
			WithAttachment("event.ics", "text/calendar", []byte("BEGIN:VCALENDAR"))
		err := service.Send(msg)
		assert.NoError(t, err)
		assert.Len(t, mock.LastEmail.Attachment, 1)
	})

	t.Run("Failure", func(t *testing.T) {
		service.(*BrevoService).client = &MockBrevoClient{SendError: fmt.Errorf("failed to send email")}

		err := service.Send(NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content"))
		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
	})
}

// TestBrevoService_SendWithResult tests the SendWithResult method of BrevoService
func TestBrevoService_SendWithResult(t *testing.T) {
	service := NewBrevoService("test_api_key", "test_sender_name", "test_sender_email")

	t.Run("Success - returns message ID verbatim", func(t *testing.T) {
		messageID := "<202406301200.123456789@smtp-relay.mailin.fr>"
		service.(*BrevoService).client = &MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{MessageId: messageID},
		}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content")
		result, err := service.(*BrevoService).SendWithResult(msg)

		assert.NoError(t, err)
		assert.Equal(t, messageID, result.MessageID)
	})

	t.Run("Failure - returns empty result and error", func(t *testing.T) {
		service.(*BrevoService).client = &MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{MessageId: "<should-be-ignored@mailin.fr>"},
			SendError:    fmt.Errorf("failed to send email"),
		}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content")
		result, err := service.(*BrevoService).SendWithResult(msg)

		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
		assert.Equal(t, SendResult{}, result)
	})
}
