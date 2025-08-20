package goat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"another.test@domain.co", true},
		{"@missingusername.com", false},
		{"missingatsign.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsEmailValid(tt.email))
		})
	}
}
