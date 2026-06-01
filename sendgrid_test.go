package goat

import (
	"context"
	"fmt"
	"testing"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

type MockSendgridClient struct {
	SendResponse *rest.Response
	SendError    error
}

func (m MockSendgridClient) Send(email *mail.SGMailV3) (*rest.Response, error) {
	return m.SendResponse, m.SendError
}

func (m MockSendgridClient) SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error) {
	return m.SendResponse, m.SendError
}

// TestNewSendgridService tests the NewSendgridService function
func TestNewSendgridService(t *testing.T) {
	service := NewSendgridService("test_api_key", "test_sender_name", "test_sender_email")

	assert.NotNil(t, service)
	assert.IsType(t, &SendgridService{}, service)
}

// TestSendgridService_Send tests the Send method of SendgridService
func TestSendgridService_Send(t *testing.T) {
	service := NewSendgridService("test_api_key", "test_sender_name", "test_sender_email")

	// Test case where sending email succeeds
	t.Run("Success", func(t *testing.T) {
		service.(*SendgridService).client = MockSendgridClient{
			SendResponse: &rest.Response{},
			SendError:    nil,
		}

		err := service.Send("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content")
		assert.NoError(t, err)
	})

	// Test case where sending email fails
	t.Run("Failure", func(t *testing.T) {
		service.(*SendgridService).client = MockSendgridClient{
			SendResponse: &rest.Response{},
			SendError:    fmt.Errorf("failed to send email"),
		}

		err := service.Send("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content")
		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
	})
}

// TestSendgridService_SendWithReplyTo tests the SendWithReplyTo method of SendgridService
func TestSendgridService_SendWithReplyTo(t *testing.T) {
	service := NewSendgridService("test_api_key", "test_sender_name", "test_sender_email")

	t.Run("Success without reply-to address", func(t *testing.T) {
		service.(*SendgridService).client = MockSendgridClient{
			SendResponse: &rest.Response{},
			SendError:    nil,
		}

		err := service.SendWithReplyTo("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content", ReplyTo{})
		assert.NoError(t, err)
	})

	t.Run("Success with reply-to address", func(t *testing.T) {
		service.(*SendgridService).client = MockSendgridClient{
			SendResponse: &rest.Response{},
			SendError:    nil,
		}

		err := service.SendWithReplyTo("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content", ReplyTo{
			Name:    "Reply Name",
			Address: "reply@example.com",
		})
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		service.(*SendgridService).client = MockSendgridClient{
			SendResponse: &rest.Response{},
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
