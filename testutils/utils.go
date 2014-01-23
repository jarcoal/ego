// random test data from http://www.databasetestdata.com/
// hopefully it doesn't collide with any actual addresses :)
package testutils

import (
	"github.com/jarcoal/ego"
	"net/mail"
	"os"
	"testing"
)

func TestRecipients() []*ego.Recipient {
	return []*ego.Recipient{
		{
			Email:           &mail.Address{"Sandy Schulist", "zane@anastacio.co.uk"},
			TemplateContext: map[string]string{"name": "Sandy"},
		},
		{
			Email:           &mail.Address{"Rocio Christiansen", "retta.ankunding@fletcher.biz"},
			TemplateContext: map[string]string{"name": "Rocio"},
		},
		{
			Email:           &mail.Address{"Abigale Gleason", "freida@orpha.info"},
			TemplateContext: map[string]string{"name": "Abigale"},
		},
		{
			Email:           &mail.Address{"Garland Spencer", "corrine@remington.io"},
			TemplateContext: map[string]string{"name": "Garland"},
		},
		{
			Email:           &mail.Address{"Tad Will", "ludwig@paula.co.uk"},
			TemplateContext: map[string]string{"name": "Tad"},
		},
		{
			Email:           &mail.Address{"Chad Ritchie", "kathryne_ankunding@uriel.biz"},
			TemplateContext: map[string]string{"name": "Chad"},
		},
		{
			Email:           &mail.Address{"Junius Boehm", "baylee.fadel@ellis.info"},
			TemplateContext: map[string]string{"name": "Junius"},
		},
		{
			Email:           &mail.Address{"Edison Kris", "alyce_rutherford@gennaro.biz"},
			TemplateContext: map[string]string{"name": "Edison"},
		},
	}
}

func TestEmail() *ego.Email {
	e := ego.NewEmail()
	e.To = TestRecipients()
	e.From = &mail.Address{Name: "Nyasia Block", Address: "jade@austen.name"}
	e.ReplyTo = &mail.Address{Name: "Erin Dare", Address: "otilia@hermina.io"}
	e.Subject = "Test Subject"
	e.HtmlBody = "<h1>Test Body</h1>"
	e.TextBody = "Test Body"
	e.Tags = []string{"really", "important", "message"}

	return e
}

func TestAttachment(t *testing.T) *ego.Attachment {
	file, err := os.Open("../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	return &ego.Attachment{"test-file.txt", "text/plain", file}
}
