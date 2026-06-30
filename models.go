package goat

// SendResult holds metadata returned by a provider after sending an email.
type SendResult struct {
	// MessageID is the provider message identifier, kept verbatim
	// (e.g. Brevo returns it wrapped in chevrons: "<xxx@smtp-relay.mailin.fr>").
	MessageID string
}

// ReplyTo holds the reply-to name and address for an email.
type ReplyTo struct {
	Name    string
	Address string
}

// Attachment represents a file attached to an email.
// Content holds the raw (un-encoded) bytes; the library base64-encodes it per provider.
// Set ContentID to embed the attachment inline (referenced from HTML as cid:<ContentID>).
// Note: Brevo (brevo-go v1.1.3) ignores ContentType and ContentID — it infers the MIME type
// from Filename and cannot embed inline; such attachments are sent as regular attachments.
type Attachment struct {
	Filename    string
	ContentType string // MIME type, e.g. "application/pdf", "text/calendar"
	Content     []byte // raw bytes (not base64)
	ContentID   string // optional; when set, the attachment is inline
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
	Attachments      []Attachment
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

// WithAttachment adds a standard file attachment and returns the message for chaining.
func (m *EmailMessage) WithAttachment(filename, contentType string, content []byte) *EmailMessage {
	m.Attachments = append(m.Attachments, Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
	})
	return m
}

// WithInlineAttachment adds an inline attachment referenced from HTML as cid:<contentID>
// and returns the message for chaining.
func (m *EmailMessage) WithInlineAttachment(filename, contentType string, content []byte, contentID string) *EmailMessage {
	m.Attachments = append(m.Attachments, Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
		ContentID:   contentID,
	})
	return m
}

// WithAttachments appends pre-built attachments and returns the message for chaining.
func (m *EmailMessage) WithAttachments(attachments ...Attachment) *EmailMessage {
	m.Attachments = append(m.Attachments, attachments...)
	return m
}
