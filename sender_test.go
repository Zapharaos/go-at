package goat

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Manual mock implementation for SenderService
type MockSenderService struct {
	SendFunc  func(emailTo, subject, plainTextContent, htmlContent string) error
	SendCalls []SendCall
	mu        sync.Mutex
}

type SendCall struct {
	EmailTo          string
	Subject          string
	PlainTextContent string
	HTMLContent      string
}

func NewMockSenderService() *MockSenderService {
	return &MockSenderService{
		SendCalls: make([]SendCall, 0),
	}
}

func (m *MockSenderService) Send(emailTo, subject, plainTextContent, htmlContent string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendCalls = append(m.SendCalls, SendCall{
		EmailTo:          emailTo,
		Subject:          subject,
		PlainTextContent: plainTextContent,
		HTMLContent:      htmlContent,
	})

	if m.SendFunc != nil {
		return m.SendFunc(emailTo, subject, plainTextContent, htmlContent)
	}
	return nil
}

func (m *MockSenderService) GetSendCalls() []SendCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]SendCall{}, m.SendCalls...)
}

func (m *MockSenderService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SendCalls = m.SendCalls[:0]
	m.SendFunc = nil
}

// TestSetSenderService tests the SetSenderService function
// It verifies that the global service can be replaced and restored correctly.
func TestSetSenderService(t *testing.T) {
	t.Run("Basic replacement and restoration", func(t *testing.T) {
		// Ensure we start with a clean state by setting a known service
		cleanupService := NewMockSenderService()
		initialRestore := SetSenderService(cleanupService)
		defer initialRestore() // Restore whatever was there originally at the end

		// Mock
		m := NewMockSenderService()

		// Store the current service (our cleanup service)
		originalService := GetSenderService()

		// Replace the global service with a mock service
		restore := SetSenderService(m)

		// Ensure the global service is replaced
		assert.Equal(t, m, GetSenderService())

		// Restore the previous global service
		restore()
		assert.Equal(t, originalService, GetSenderService())
		// The restoration worked correctly if we get back to the original service
	})

	t.Run("Multiple replacements work correctly", func(t *testing.T) {
		// Ensure we start with a clean state
		cleanupService := NewMockSenderService()
		initialRestore := SetSenderService(cleanupService)
		defer initialRestore()

		mock1 := NewMockSenderService()
		mock2 := NewMockSenderService()

		// Store original service
		originalService := GetSenderService()

		// Set first service
		restore1 := SetSenderService(mock1)
		assert.Equal(t, mock1, GetSenderService())

		// Set second service
		restore2 := SetSenderService(mock2)
		assert.Equal(t, mock2, GetSenderService())

		// Restore to first service
		restore2()
		assert.Equal(t, mock1, GetSenderService())

		// Restore to original
		restore1()
		assert.Equal(t, originalService, GetSenderService())
	})

	t.Run("Restore function can be called multiple times safely", func(t *testing.T) {
		// Ensure we start with a clean state
		cleanupService := NewMockSenderService()
		initialRestore := SetSenderService(cleanupService)
		defer initialRestore()

		mock := NewMockSenderService()

		originalService := GetSenderService()
		restore := SetSenderService(mock)

		assert.Equal(t, mock, GetSenderService())

		// Call restore multiple times
		restore()
		assert.Equal(t, originalService, GetSenderService())

		restore() // Should not panic or cause issues
		assert.Equal(t, originalService, GetSenderService())
	})
}

// TestGetSenderService tests the GetSenderService function
// It verifies that the global service can be accessed correctly.
func TestGetSenderService(t *testing.T) {
	t.Run("Returns the current global service", func(t *testing.T) {
		// Mock
		m := NewMockSenderService()

		// Replace the global service with a mock service
		restore := SetSenderService(m)
		defer restore()

		// Access the global service
		service := GetSenderService()
		assert.Equal(t, m, service)
	})

	t.Run("Returns nil when no service is set", func(t *testing.T) {
		// Temporarily clear the global service
		restore := SetSenderService(nil)
		defer restore()

		service := GetSenderService()
		assert.Nil(t, service)
	})
}

// TestGlobalFunctions tests the global wrapper functions
func TestGlobalFunctions(t *testing.T) {
	mockService := NewMockSenderService()

	restore := SetSenderService(mockService)
	defer restore()

	t.Run("Send calls service Send", func(t *testing.T) {
		expectedErr := assert.AnError

		// Set up the mock to return the expected error
		mockService.SendFunc = func(emailTo, subject, plainTextContent, htmlContent string) error {
			return expectedErr
		}

		err := Send("", "", "", "")

		assert.Equal(t, expectedErr, err)

		// Verify the call was made
		calls := mockService.GetSendCalls()
		assert.Len(t, calls, 1)
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

		// All reads should return the same service
		for i, result := range results {
			assert.Equal(t, mockService, result, "Reader %d got unexpected result", i)
		}
	})

	t.Run("Concurrent writes are safe", func(t *testing.T) {
		const numWriters = 10
		var wg sync.WaitGroup
		mocks := make([]*MockSenderService, numWriters)
		restoreFuncs := make([]func(), numWriters)

		// Create mock services
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

		// The final service should be one of the mocks
		finalService := GetSenderService()
		found := false
		for _, mock := range mocks {
			if finalService == mock {
				found = true
				break
			}
		}
		assert.True(t, found, "Final service should be one of the set mocks")

		// Clean up (call restore functions)
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

		wg.Add(numOperations * 2) // readers + writers

		// Start readers
		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				service := GetSenderService()
				assert.NotNil(t, service) // Should always get some service
			}()
		}

		// Start writers
		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				newMock := NewMockSenderService()
				restore := SetSenderService(newMock)
				// Immediately restore to avoid leaving test in inconsistent state
				restore()
			}()
		}

		wg.Wait()
		// Test should complete without data races or panics
	})
}

// TestGlobalFunctionsWithNilService tests behavior when no service is set
func TestGlobalFunctionsWithNilService(t *testing.T) {
	// Temporarily clear the global service
	restore := SetSenderService(nil)
	defer restore()

	t.Run("Send panics with nil service", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = Send("", "", "", "")
		})
	})
}
