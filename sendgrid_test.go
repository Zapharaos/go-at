package goat

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

type MockSendgridClient struct {
	SendResponse *rest.Response
	SendError    error
	LastEmail    *mail.SGMailV3
}

func (m *MockSendgridClient) Send(email *mail.SGMailV3) (*rest.Response, error) {
	m.LastEmail = email
	return m.SendResponse, m.SendError
}

func (m *MockSendgridClient) SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error) {
	m.LastEmail = email
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

	t.Run("Success - plain message", func(t *testing.T) {
		service.(*SendgridService).client = &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}

		err := service.Send(NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content"))
		assert.NoError(t, err)
	})

	t.Run("Success - with reply-to", func(t *testing.T) {
		service.(*SendgridService).client = &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with headers", func(t *testing.T) {
		service.(*SendgridService).client = &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>, <https://example.com/unsubscribe>")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with reply-to and headers combined", func(t *testing.T) {
		service.(*SendgridService).client = &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>")
		err := service.Send(msg)
		assert.NoError(t, err)
	})

	t.Run("Success - with attachment", func(t *testing.T) {
		mock := &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}
		service.(*SendgridService).client = mock

		content := []byte("PDF bytes")
		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithAttachment("invoice.pdf", "application/pdf", content)
		err := service.Send(msg)
		assert.NoError(t, err)

		assert.Len(t, mock.LastEmail.Attachments, 1)
		att := mock.LastEmail.Attachments[0]
		assert.Equal(t, "invoice.pdf", att.Filename)
		assert.Equal(t, "application/pdf", att.Type)
		assert.Equal(t, "attachment", att.Disposition)
		assert.Equal(t, base64.StdEncoding.EncodeToString(content), att.Content)
		assert.Empty(t, att.ContentID)
	})

	t.Run("Success - with inline attachment", func(t *testing.T) {
		mock := &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}
		service.(*SendgridService).client = mock

		content := []byte("PNG bytes")
		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "<img src=\"cid:logo\">").
			WithInlineAttachment("logo.png", "image/png", content, "logo")
		err := service.Send(msg)
		assert.NoError(t, err)

		assert.Len(t, mock.LastEmail.Attachments, 1)
		att := mock.LastEmail.Attachments[0]
		assert.Equal(t, "logo.png", att.Filename)
		assert.Equal(t, "inline", att.Disposition)
		assert.Equal(t, "logo", att.ContentID)
	})

	t.Run("Success - attachment with reply-to and headers combined", func(t *testing.T) {
		mock := &MockSendgridClient{SendResponse: &rest.Response{}, SendError: nil}
		service.(*SendgridService).client = mock

		msg := NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content").
			WithReplyTo("Reply Name", "reply@example.com").
			WithHeader("List-Unsubscribe", "<mailto:unsubscribe@example.com>").
			WithAttachment("event.ics", "text/calendar", []byte("BEGIN:VCALENDAR"))
		err := service.Send(msg)
		assert.NoError(t, err)
		assert.Len(t, mock.LastEmail.Attachments, 1)
	})

	t.Run("Failure", func(t *testing.T) {
		service.(*SendgridService).client = &MockSendgridClient{SendResponse: &rest.Response{}, SendError: fmt.Errorf("failed to send email")}

		err := service.Send(NewEmailMessage("test@example.com", "Test Subject", "Test Plain Text", "Test HTML Content"))
		assert.Error(t, err)
		assert.Equal(t, "failed to send email", err.Error())
	})
}
