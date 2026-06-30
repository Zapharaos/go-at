package goat

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Manual mock implementation for SenderService
type MockSenderService struct {
	SendFunc           func(message *EmailMessage) error
	SendWithResultFunc func(message *EmailMessage) (SendResult, error)
	SendCalls          []*EmailMessage
	mu                 sync.Mutex
}

func NewMockSenderService() *MockSenderService {
	return &MockSenderService{
		SendCalls: make([]*EmailMessage, 0),
	}
}

func (m *MockSenderService) Send(message *EmailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendCalls = append(m.SendCalls, message)

	if m.SendFunc != nil {
		return m.SendFunc(message)
	}
	return nil
}

func (m *MockSenderService) SendWithResult(message *EmailMessage) (SendResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendCalls = append(m.SendCalls, message)

	if m.SendWithResultFunc != nil {
		return m.SendWithResultFunc(message)
	}
	return SendResult{}, nil
}

func (m *MockSenderService) GetSendCalls() []*EmailMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]*EmailMessage{}, m.SendCalls...)
}

func (m *MockSenderService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SendCalls = m.SendCalls[:0]
	m.SendFunc = nil
	m.SendWithResultFunc = nil
}

// TestSetSenderService tests the SetSenderService function
// It verifies that the global service can be replaced and restored correctly.
func TestSetSenderService(t *testing.T) {
	t.Run("Basic replacement and restoration", func(t *testing.T) {
		cleanupService := NewMockSenderService()
		initialRestore := SetSenderService(cleanupService)
		defer initialRestore()

		m := NewMockSenderService()
		originalService := GetSenderService()

		restore := SetSenderService(m)
		assert.Equal(t, m, GetSenderService())

		restore()
		assert.Equal(t, originalService, GetSenderService())
	})

	t.Run("Multiple replacements work correctly", func(t *testing.T) {
		cleanupService := NewMockSenderService()
		initialRestore := SetSenderService(cleanupService)
		defer initialRestore()

		mock1 := NewMockSenderService()
		mock2 := NewMockSenderService()
		originalService := GetSenderService()

		restore1 := SetSenderService(mock1)
		assert.Equal(t, mock1, GetSenderService())

		restore2 := SetSenderService(mock2)
		assert.Equal(t, mock2, GetSenderService())

		restore2()
		assert.Equal(t, mock1, GetSenderService())

		restore1()
		assert.Equal(t, originalService, GetSenderService())
	})

	t.Run("Restore function can be called multiple times safely", func(t *testing.T) {
		cleanupService := NewMockSenderService()
		initialRestore := SetSenderService(cleanupService)
		defer initialRestore()

		mock := NewMockSenderService()
		originalService := GetSenderService()
		restore := SetSenderService(mock)

		assert.Equal(t, mock, GetSenderService())

		restore()
		assert.Equal(t, originalService, GetSenderService())

		restore() // Should not panic or cause issues
		assert.Equal(t, originalService, GetSenderService())
	})
}

// TestGetSenderService tests the GetSenderService function
func TestGetSenderService(t *testing.T) {
	t.Run("Returns the current global service", func(t *testing.T) {
		m := NewMockSenderService()
		restore := SetSenderService(m)
		defer restore()

		assert.Equal(t, m, GetSenderService())
	})

	t.Run("Returns nil when no service is set", func(t *testing.T) {
		restore := SetSenderService(nil)
		defer restore()

		assert.Nil(t, GetSenderService())
	})
}

// TestGlobalFunctions tests the global wrapper functions
func TestGlobalFunctions(t *testing.T) {
	mockService := NewMockSenderService()

	restore := SetSenderService(mockService)
	defer restore()

	t.Run("Send calls service Send", func(t *testing.T) {
		mockService.Reset()
		expectedErr := assert.AnError

		mockService.SendFunc = func(message *EmailMessage) error {
			return expectedErr
		}

		msg := NewEmailMessage("to@example.com", "Subject", "plain", "<b>html</b>")
		err := Send(msg)

		assert.Equal(t, expectedErr, err)
		assert.Len(t, mockService.GetSendCalls(), 1)
	})

	t.Run("SendWithResult calls service SendWithResult", func(t *testing.T) {
		mockService.Reset()
		expected := SendResult{MessageID: "<msg-123@smtp-relay.mailin.fr>"}

		mockService.SendWithResultFunc = func(message *EmailMessage) (SendResult, error) {
			return expected, nil
		}

		msg := NewEmailMessage("to@example.com", "Subject", "plain", "<b>html</b>")
		result, err := SendWithResult(msg)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Len(t, mockService.GetSendCalls(), 1)
	})

	t.Run("Send passes message fields correctly", func(t *testing.T) {
		mockService.Reset()

		headers := map[string]string{"List-Unsubscribe": "<mailto:unsub@example.com>"}
		content := []byte("PDF bytes")
		msg := NewEmailMessage("to@example.com", "Subject", "plain", "<b>html</b>").
			WithReplyTo("Reply Name", "reply@example.com").
			WithHeaders(headers).
			WithAttachment("invoice.pdf", "application/pdf", content)

		_ = Send(msg)

		calls := mockService.GetSendCalls()
		assert.Len(t, calls, 1)
		assert.Equal(t, "to@example.com", calls[0].To)
		assert.Equal(t, "Subject", calls[0].Subject)
		assert.Equal(t, "plain", calls[0].PlainTextContent)
		assert.Equal(t, "<b>html</b>", calls[0].HTMLContent)
		assert.Equal(t, &ReplyTo{Name: "Reply Name", Address: "reply@example.com"}, calls[0].ReplyTo)
		assert.Equal(t, headers, calls[0].Headers)
		assert.Equal(t, []Attachment{{Filename: "invoice.pdf", ContentType: "application/pdf", Content: content}}, calls[0].Attachments)
	})
}

// TestConcurrentAccess tests thread safety of the global service
func TestConcurrentAccess(t *testing.T) {
	t.Run("Concurrent reads are safe", func(t *testing.T) {
		mockService := NewMockSenderService()
		restore := SetSenderService(mockService)
		defer restore()

		const numReaders = 100
		var wg sync.WaitGroup
		results := make([]SenderService, numReaders)

		wg.Add(numReaders)
		for i := 0; i < numReaders; i++ {
			go func(idx int) {
				defer wg.Done()
				results[idx] = GetSenderService()
			}(i)
		}
		wg.Wait()

		for i, result := range results {
			assert.Equal(t, mockService, result, "Reader %d got unexpected result", i)
		}
	})

	t.Run("Concurrent writes are safe", func(t *testing.T) {
		const numWriters = 10
		var wg sync.WaitGroup
		mocks := make([]*MockSenderService, numWriters)
		restoreFuncs := make([]func(), numWriters)

		for i := 0; i < numWriters; i++ {
			mocks[i] = NewMockSenderService()
		}

		wg.Add(numWriters)
		for i := 0; i < numWriters; i++ {
			go func(idx int) {
				defer wg.Done()
				restoreFuncs[idx] = SetSenderService(mocks[idx])
			}(i)
		}
		wg.Wait()

		finalService := GetSenderService()
		found := false
		for _, mock := range mocks {
			if finalService == mock {
				found = true
				break
			}
		}
		assert.True(t, found, "Final service should be one of the set mocks")

		for _, restore := range restoreFuncs {
			if restore != nil {
				restore()
			}
		}
	})

	t.Run("Mixed concurrent reads and writes are safe", func(t *testing.T) {
		mockService := NewMockSenderService()
		const numOperations = 50
		var wg sync.WaitGroup

		originalRestore := SetSenderService(mockService)
		defer originalRestore()

		wg.Add(numOperations * 2)

		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				assert.NotNil(t, GetSenderService())
			}()
		}

		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				newMock := NewMockSenderService()
				restore := SetSenderService(newMock)
				restore()
			}()
		}

		wg.Wait()
	})
}

// TestGlobalFunctionsWithNilService tests behavior when no service is set
func TestGlobalFunctionsWithNilService(t *testing.T) {
	restore := SetSenderService(nil)
	defer restore()

	t.Run("Send panics with nil service", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = Send(NewEmailMessage("", "", "", ""))
		})
	})
}
