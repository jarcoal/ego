package ego

import (
	"testing"
)

func TestEmailAddAttachment(t *testing.T) {
	e := NewEmail()

	if len(e.Attachments) != 0 {
		t.FailNow()
	}

	e.AddAttachment("test-attachment", "text/plain", nil)

	if len(e.Attachments) != 1 {
		t.FailNow()
	}
}

func TestEmailAddRecipient(t *testing.T) {
	e := NewEmail()

	if len(e.To) != 0 {
		t.FailNow()
	}

	e.AddRecipient("test recipient", "test@test.com")

	if len(e.To) != 1 {
		t.FailNow()
	}
}
