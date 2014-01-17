// Dummy Backend
// Doesn't do anything except optionally log out the emails it receives.
package dummy

import (
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/backends"
	"strings"
)

var _ backends.Backend = (*dummyBackend)(nil)

// to have the dummyBackend log given emails, it needs to be instantiated with this function type
type logger func(format string, vars ...interface{})

func NewBackend(l logger) backends.Backend {
	return &dummyBackend{l}
}

type dummyBackend struct {
	log logger
}

func (d *dummyBackend) DispatchEmail(e *ego.Email) error {
	if d.log == nil {
		return nil
	}

	recipients := []string{}

	for _, recip := range e.To {
		recipients = append(recipients, recip.String())
	}

	d.log("To: %s", strings.Join(recipients, ", "))
	d.log("From: %s", e.From)
	d.log("Subject: %s", e.Subject)
	d.log("TrackClicks: %t", e.TrackClicks)
	d.log("TrackOpens: %t", e.TrackOpens)
	d.log("VisibleRecipients: %t", e.VisibleRecipients)

	if len(e.Tags) > 0 {
		d.log("Tags: %s", strings.Join(e.Tags, ", "))
	}

	if e.TemplateId != "" {
		d.log("TemplateId: %s", e.TemplateId)
		d.log("TemplateContext: %s", e.TemplateContext)
	}

	if e.SubAccount != "" {
		d.log("SubAccount: %s", e.SubAccount)
	}

	if !e.DeliveryTime.IsZero() {
		d.log("DeliveryTime: %s", e.DeliveryTime)
	}

	if len(e.Attachments) > 0 {
		attachmentNames := []string{}

		for _, attachment := range e.Attachments {
			attachmentNames = append(attachmentNames, attachment.Name)
		}

		d.log("Attachments: %s", strings.Join(attachmentNames, ", "))
	}

	d.log("TextBody: %s", e.TextBody)
	d.log("HtmlBody: %s", e.HtmlBody)

	return nil
}
