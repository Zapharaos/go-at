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

	// Test case where sending email succeeds
	t.Run("Success", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{},
			SendError:    nil,
		}

		err := service.Send("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content")
		assert.NoError(t, err)
	})

	// Test case where sending email fails
	t.Run("Failure", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{},
			SendError:    fmt.Errorf("failed to send email"),
		}

		err := service.Send("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content")
		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
	})
}

// TestBrevoService_SendWithReplyTo tests the SendWithReplyTo method of BrevoService
func TestBrevoService_SendWithReplyTo(t *testing.T) {
	service := NewBrevoService("test_api_key", "test_sender_name", "test_sender_email")

	t.Run("Success without reply-to address", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{},
			SendError:    nil,
		}

		err := service.SendWithReplyTo("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content", ReplyTo{})
		assert.NoError(t, err)
	})

	t.Run("Success with reply-to address", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{},
			SendError:    nil,
		}

		err := service.SendWithReplyTo("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content", ReplyTo{
			Name:    "Reply Name",
			Address: "reply@example.com",
		})
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		service.(*BrevoService).client = MockBrevoClient{
			SendResponse: brevo.CreateSmtpEmail{},
			SendError:    fmt.Errorf("failed to send email"),
		}

		err := service.SendWithReplyTo("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content", ReplyTo{
			Name:    "Reply Name",
			Address: "reply@example.com",
		})
		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
	})
}
