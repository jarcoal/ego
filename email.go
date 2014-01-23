package ego

import (
	"io"
	"net/mail"
	"time"
)

func NewEmail() *Email {
	return &Email{
		To:              make([]*Recipient, 0),
		Cc:              make([]*Recipient, 0),
		Bcc:             make([]*Recipient, 0),
		TrackClicks:     true,
		TrackOpens:      true,
		Attachments:     make([]*Attachment, 0),
		TemplateContext: make(map[string]string),
	}
}

type Email struct {
	// Basic Email Sender/Receiver Information
	To, Cc, Bcc   []*Recipient
	From, ReplyTo *mail.Address

	// Email Subject
	Subject string

	// The separate chunks
	HtmlBody, TextBody string

	// Files/data to be attached to the email
	Attachments []*Attachment

	// Many email services offer a tagging system for emails so
	// they can be grouped for analytics.
	Tags []string

	// Many email services offer a templating system where you can pass and ID
	// and a dictionary of context and they will render the template on their end.
	TemplateId      string
	TemplateContext map[string]string

	// Sender-side user identifier.  Many email services will let you assign identifiers for users
	// so that their send volume can be monitored on an individual level.
	SubAccount string

	// Behavior tracking preferences
	TrackClicks, TrackOpens bool

	// Many services allow you to queue email to be delivered at a specific date.
	DeliveryTime time.Time

	// Can recipients see the names of other recipients in the "from" header?
	// This almost always should be `false`.
	VisibleRecipients bool
}

// AddAttachment is a convenience method for adding attachments to the message
func (e *Email) AddAttachment(name, mimetype string, data io.Reader) {
	e.Attachments = append(e.Attachments, &Attachment{name, mimetype, data})
}

// AddRecipient is a convenience method for adding recipients to the message
func (e *Email) AddRecipient(name, email string, tmplCtx map[string]string) {
	e.To = append(e.To, &Recipient{
		Email: &mail.Address{
			Name:    name,
			Address: email,
		},
		TemplateContext: tmplCtx,
	})
}

// Attachment represents a piece of data to be attached to an email.
type Attachment struct {
	Name, Mimetype string
	Data           io.Reader
}

// Recipient represents a single recipient in an email.
type Recipient struct {
	Email           *mail.Address
	TemplateContext map[string]string
}
