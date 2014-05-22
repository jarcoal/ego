// Test helpers and dummy data.

package testutils

import (
	"github.com/jarcoal/ego"
	"net/mail"
	"os"
	"testing"
)

// TestRecipients returns a batch of dummy names/emails to be used in tests.
func TestRecipients() []*ego.Recipient {
	return []*ego.Recipient{
		{
			Email:           &mail.Address{Name: "Sandy Schulist", Address: "zane@anastacio.co.uk"},
			TemplateContext: map[string]string{"name": "Sandy"},
		},
		{
			Email:           &mail.Address{Name: "Rocio Christiansen", Address: "retta.ankunding@fletcher.biz"},
			TemplateContext: map[string]string{"name": "Rocio"},
		},
		{
			Email:           &mail.Address{Name: "Abigale Gleason", Address: "freida@orpha.info"},
			TemplateContext: map[string]string{"name": "Abigale"},
		},
		{
			Email:           &mail.Address{Name: "Garland Spencer", Address: "corrine@remington.io"},
			TemplateContext: map[string]string{"name": "Garland"},
		},
		{
			Email:           &mail.Address{Name: "Tad Will", Address: "ludwig@paula.co.uk"},
			TemplateContext: map[string]string{"name": "Tad"},
		},
		{
			Email:           &mail.Address{Name: "Chad Ritchie", Address: "kathryne_ankunding@uriel.biz"},
			TemplateContext: map[string]string{"name": "Chad"},
		},
		{
			Email:           &mail.Address{Name: "Junius Boehm", Address: "baylee.fadel@ellis.info"},
			TemplateContext: map[string]string{"name": "Junius"},
		},
		{
			Email:           &mail.Address{Name: "Edison Kris", Address: "alyce_rutherford@gennaro.biz"},
			TemplateContext: map[string]string{"name": "Edison"},
		},
	}
}

// TestEmail returns an Email populated with dummy data to be used in tests.
func TestEmail() *ego.Email {
	e := ego.NewEmail()
	e.To = TestRecipients()
	e.From = &mail.Address{Name: "Nyasia Block", Address: "jade@austen.name"}
	e.ReplyTo = &mail.Address{Name: "Erin Dare", Address: "otilia@hermina.io"}
	e.Subject = "Test Subject"
	e.HTMLBody = "<h1>Test Body</h1>"
	e.TextBody = "Test Body"
	e.Tags = []string{"really", "important", "message"}

	return e
}

// TestAttachment creates an attachment from a file.
func TestAttachment(t *testing.T) *ego.Attachment {
	file, err := os.Open("../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	return &ego.Attachment{Name: "test-file.txt", Mimetype: "text/plain", Data: file}
}
