package mail

import (
	"context"
	"errors"
	"time"
)

// DefaultCharset is a default character encoding for the email.
const DefaultCharset string = "UTF-8"

var (
	// ErrMessageNil returned when message is nil.
	ErrMessageNil = errors.New("mail: message is nil")
	// ErrEmptySender returned when sender is empty.
	ErrEmptySender = errors.New("mail: empty sender")
	// ErrEmptyRecipient returned when recipient is empty.
	ErrEmptyRecipient = errors.New("mail: empty recipient")
	// ErrEmptySubject returned when subject is empty.
	ErrEmptySubject = errors.New("mail: empty subject")
	// ErrEmptyHTMLBody returned when html body is empty.
	ErrEmptyHTMLBody = errors.New("mail: empty html body")
	// ErrEmptyTextBody returned when text body is empty.
	ErrEmptyTextBody = errors.New("mail: empty text body")
)

// Message holds email details.
type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	HTMLBody  string `json:"htmlbody"`
	TextBody  string `json:"textbody"`
	Charset   string `json:"charset"`
}

// NewMessage returns a new instance of `Message`.
func NewMessage(sender, recipient, subject, htmlBody, textBody, charset string) *Message {
	m := &Message{
		Sender:    sender,
		Recipient: recipient,
		Subject:   subject,
		HTMLBody:  htmlBody,
		TextBody:  textBody,
		Charset:   charset,
	}
	if m.Charset == "" {
		m.Charset = DefaultCharset
	}
	return m
}

func empty(s string) bool { return s == "" }

// Valid checks mail fields.
func (m *Message) Valid() error {
	if m == nil {
		return ErrMessageNil
	}
	if empty(m.Sender) {
		return ErrEmptySender
	}
	if empty(m.Recipient) {
		return ErrEmptyRecipient
	}
	if empty(m.Subject) {
		return ErrEmptySubject
	}
	if empty(m.HTMLBody) {
		return ErrEmptyHTMLBody
	}
	if empty(m.TextBody) {
		return ErrEmptyTextBody
	}
	return nil
}

// Mailer implements mail send.
type Mailer interface {
	// SendMail ...
	// TODO: description
	SendMail(ctx context.Context, message *Message) (*time.Time, error)
}
