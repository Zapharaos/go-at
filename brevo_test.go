package goat

import (
	"context"
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
}

func (m MockBrevoClient) SendTransacEmail(ctx context.Context, sendSmtpEmail brevo.SendSmtpEmail) (brevo.CreateSmtpEmail, *http.Response, error) {
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
		service.(*BrevoService).client = MockBrevoClient{SendError: nil}

		err := service.Send(NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content"))
		assert.NoError(t, err)
	})

	t.Run("Success - with reply-to", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with headers", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>, <https://example.com/unsubscribe>")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with reply-to and headers combined", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{SendError: fmt.Errorf("failed to send email")}

		err := service.Send(NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content"))
		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
	})
}
