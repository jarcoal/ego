package mandrill

import (
	"encoding/base64"
	"github.com/jarcoal/ego/testutils"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

var b = mandrillBackend{"abc123"}

func TestWrapper(t *testing.T) {
	e := testutils.TestEmail()

	// set some additional properties
	e.DeliveryTime = time.Now()
	e.TemplateId = "test-template"

	wrapper, err := b.mandrillWrapperForEmail(e)
	if err != nil {
		t.FailNow()
	}

	if wrapper.SendAt != e.DeliveryTime.Format(MANDRILL_DELIVERY_TIME_FMT) {
		t.FailNow()
	}

	if wrapper.TemplateName != e.TemplateId {
		t.FailNow()
	}
}

func TestEmail(t *testing.T) {
	e := testutils.TestEmail()

	wrapper, err := b.mandrillWrapperForEmail(e)
	if err != nil {
		t.FailNow()
	}

	me := wrapper.Message

	if e.Subject != me.Subject {
		t.FailNow()
	}

	if e.HtmlBody != me.Html {
		t.FailNow()
	}

	if e.TextBody != me.Text {
		t.FailNow()
	}

	if e.From.Name != me.FromName {
		t.FailNow()
	}

	if e.From.Address != me.FromEmail {
		t.FailNow()
	}

	if e.ReplyTo.Address != me.Headers["Reply-To"] {
		t.FailNow()
	}

	if e.TrackClicks != me.TrackClicks {
		t.FailNow()
	}

	if e.TrackOpens != me.TrackOpens {
		t.FailNow()
	}

	if e.SubAccount != me.Subaccount {
		t.FailNow()
	}

	if len(e.Tags) != len(me.Tags) {
		t.FailNow()
	}

	if len(e.To) != len(me.To) {
		t.FailNow()
	}
}

func TestTemplating(t *testing.T) {
	// TestEmail adds recipients that have a piece of recipient-specific context
	e := testutils.TestEmail()

	// set some global template context
	e.TemplateContext["hello"] = "world"

	wrapper, err := b.mandrillWrapperForEmail(e)
	if err != nil {
		t.FailNow()
	}

	me := wrapper.Message

	if len(me.GlobalMergeVars) != 1 {
		t.FailNow()
	}

	if me.GlobalMergeVars[0].Name != "hello" {
		t.FailNow()
	}

	if me.GlobalMergeVars[0].Content != "world" {
		t.FailNow()
	}

	if len(me.MergeVars) != len(me.To) {
		t.FailNow()
	}
}

func TestAttachments(t *testing.T) {
	e := testutils.TestEmail()

	attachment := testutils.TestAttachment(t)
	e.Attachments = append(e.Attachments, attachment)

	wrapper, err := b.mandrillWrapperForEmail(e)
	if err != nil {
		t.FailNow()
	}

	me := wrapper.Message

	if len(me.Attachments) != len(e.Attachments) {
		t.FailNow()
	}

	mAttachment := me.Attachments[0]

	if attachment.Name != mAttachment.Name {
		t.FailNow()
	}

	if attachment.Mimetype != mAttachment.Type {
		t.FailNow()
	}

	// the mandrill backend will have read the contents of the reader,
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

	if base64.StdEncoding.EncodeToString(fileData) != mAttachment.Content {
		t.FailNow()
	}
}
