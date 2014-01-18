package postageapp

import (
	"encoding/base64"
	"github.com/jarcoal/ego/testutils"
	"io"
	"io/ioutil"
	"testing"
)

const POSTAGEAPP_API_KEY = "abc123"

var b = postageAppBackend{POSTAGEAPP_API_KEY}

func TestBackendWrapper(t *testing.T) {
	e := testutils.TestEmail()

	e.TemplateId = "test-template"
	e.TemplateContext = map[string]string{"hello": "world"}

	wrapper, err := b.wrapperForEmail(e)
	if err != nil {
		t.Fatal(err)
	}

	if wrapper.ApiKey != POSTAGEAPP_API_KEY {
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

	if wrapper.Arguments.Content["text/html"] != e.HtmlBody {
		t.FailNow()
	}

	if wrapper.Arguments.Template != e.TemplateId {
		t.FailNow()
	}

	if wrapper.Arguments.Variables["hello"] != e.TemplateContext["hello"] {
		t.FailNow()
	}
}

func TestBackendAttachments(t *testing.T) {
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
