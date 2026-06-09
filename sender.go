package goat

import (
	"sync"
)

// SenderService defines the interface for handling emails
type SenderService interface {
	Send(message *EmailMessage) error
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

// Send directly exposes the current sender service Send function.
func Send(message *EmailMessage) error {
	return GetSenderService().Send(message)
}
