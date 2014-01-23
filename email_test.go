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

	testName := "test recipient"
	testAddress := "test@test.com"
	testCtx := map[string]string{"hello": "world"}

	if len(e.To) != 0 {
		t.FailNow()
	}

	e.AddRecipient(testName, testAddress, testCtx)

	if len(e.To) != 1 {
		t.FailNow()
	}

	if e.To[0].Email.Address != testAddress {
		t.FailNow()
	}

	if e.To[0].Email.Name != testName {
		t.FailNow()
	}

	if e.To[0].TemplateContext["hello"] != "world" {
		t.FailNow()
	}
}
