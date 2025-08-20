package goat

import (
	"sync"
)

// TODO : examples + doc
// TODO : remove mocks from dep + makefile & replace with manual mock

// SenderService defines the interface for handling emails
type SenderService interface {
	Send(emailTo, subject, plainTextContent, htmlContent string) error
}

var (
	_globalServiceMu sync.RWMutex
	_globalService   SenderService
)

// SetSenderService affect a new repository to the global service singleton
func SetSenderService(service SenderService) func() {
	_globalServiceMu.Lock()
	defer _globalServiceMu.Unlock()

	prev := _globalService
	_globalService = service
	return func() { SetSenderService(prev) }
}

// GetSenderService is used to access the global service singleton
func GetSenderService() SenderService {
	_globalServiceMu.RLock()
	defer _globalServiceMu.RUnlock()

	service := _globalService
	return service
}

// Directly exposes global sender implementation

// Send Directly exposes the current sender service Send function.
func Send(emailTo, subject, plainTextContent, htmlContent string) error {
	return GetSenderService().Send(emailTo, subject, plainTextContent, htmlContent)
}
