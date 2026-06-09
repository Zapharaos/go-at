package goat

// ReplyTo holds the reply-to name and address for an email.
type ReplyTo struct {
	Name    string
	Address string
}

// EmailMessage represents an email to be sent.
// Build one with NewEmailMessage and chain With* methods for optional fields.
type EmailMessage struct {
	To               string
	Subject          string
	PlainTextContent string
	HTMLContent      string
	ReplyTo          *ReplyTo
	Headers          map[string]string
}

// NewEmailMessage creates a new EmailMessage with the required fields.
func NewEmailMessage(to, subject, plainTextContent, htmlContent string) *EmailMessage {
	return &EmailMessage{
		To:               to,
		Subject:          subject,
		PlainTextContent: plainTextContent,
		HTMLContent:      htmlContent,
	}
}

// WithReplyTo sets the reply-to address and returns the message for chaining.
func (m *EmailMessage) WithReplyTo(name, address string) *EmailMessage {
	m.ReplyTo = &ReplyTo{Name: name, Address: address}
	return m
}

// WithHeader adds a single custom header and returns the message for chaining.
func (m *EmailMessage) WithHeader(key, value string) *EmailMessage {
	if m.Headers == nil {
		m.Headers = make(map[string]string)
	}
	m.Headers[key] = value
	return m
}

// WithHeaders sets all custom headers at once and returns the message for chaining.
func (m *EmailMessage) WithHeaders(headers map[string]string) *EmailMessage {
	m.Headers = headers
	return m
}
