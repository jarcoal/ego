package postageapp

import (
	"encoding/base64"
	"github.com/jarcoal/ego/testutils"
	"io"
	"io/ioutil"
	"testing"
)

const apiKey = "abc123"

var b = postageAppBackend{apiKey}

// TestWrapper checks the JSON wrapper
func TestWrapper(t *testing.T) {
	e := testutils.TestEmail()

	wrapper, err := b.wrapperForEmail(e)
	if err != nil {
		t.Fatal(err)
	}

	if wrapper.APIKey != apiKey {
		t.FailNow()
	}

	if len(wrapper.Arguments.Recipients) != len(e.To) {
		t.FailNow()
	}

	if wrapper.Arguments.Headers["subject"] != e.Subject {
		t.FailNow()
	}

	if wrapper.Arguments.Headers["from"] != e.From.String() {
		t.FailNow()
	}

	if wrapper.Arguments.Headers["reply-to"] != e.ReplyTo.String() {
		t.FailNow()
	}

	if wrapper.Arguments.Content["text/plain"] != e.TextBody {
		t.FailNow()
	}

	if wrapper.Arguments.Content["text/html"] != e.HTMLBody {
		t.FailNow()
	}
}

// TestHeaders checks that headers are set correctly
func TestHeaders(t *testing.T) {
	e := testutils.TestEmail()

	e.Headers.Set("hello", "world")

	wrapper, err := b.wrapperForEmail(e)
	if err != nil {
		t.FailNow()
	}

	if wrapper.Arguments.Headers["hello"] != "world" {
		t.FailNow()
	}
}

// TestTemplating checks that template name/vars are set correctly
func TestTemplating(t *testing.T) {
	e := testutils.TestEmail()
	e.TemplateID = "test-template"
	e.TemplateContext = map[string]string{"hello": "world"}

	wrapper, err := b.wrapperForEmail(e)
	if err != nil {
		t.Fatal(err)
	}

	if wrapper.Arguments.Template != e.TemplateID {
		t.FailNow()
	}

	if wrapper.Arguments.Variables["hello"] != e.TemplateContext["hello"] {
		t.FailNow()
	}

	for _, ctx := range wrapper.Arguments.Recipients {
		if _, ok := ctx["name"]; !ok {
			t.FailNow()
		}
	}
}

// TestAttachments checks that attachments are added correctly
func TestAttachments(t *testing.T) {
	e := testutils.TestEmail()

	attachment := testutils.TestAttachment(t)
	e.Attachments = append(e.Attachments, attachment)

	wrapper, err := b.wrapperForEmail(e)
	if err != nil {
		t.Fatal(err)
	}

	paAttachment, ok := wrapper.Arguments.Attachments[attachment.Name]
	if !ok {
		t.FailNow()
	}

	if paAttachment.ContentType != attachment.Mimetype {
		t.FailNow()
	}

	// the postageapp backend will have read the contents of the reader,
	// so we need to seek back to the beginning
	offset, err := attachment.Data.(io.ReadSeeker).Seek(0, 0)
	if err != nil {
		t.FailNow()
	} else if offset != int64(0) {
		t.FailNow()
	}

	// check that the file data is there
	fileData, err := ioutil.ReadAll(attachment.Data)
	if err != nil {
		t.FailNow()
	}

	if paAttachment.Content != base64.StdEncoding.EncodeToString(fileData) {
		t.FailNow()
	}
}
