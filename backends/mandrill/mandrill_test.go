package backends

import (
	"encoding/base64"
	"github.com/jarcoal/ego/testutils"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

var b = mandrillBackend{"abc123"}

func TestBackendWrapper(t *testing.T) {
	e := testutils.TestEmail()

	// set some additional properties
	e.DeliveryTime = time.Now()
	e.TemplateId = "test-template"
	e.TemplateContext = map[string]string{"hello": "world"}

	me, err := b.mandrillEmailForEmail(e)
	if err != nil {
		t.FailNow()
	}

	wrapper := b.wrapMandrillEmail(e, me)

	if wrapper["send_at"] != e.DeliveryTime.Format(MANDRILL_DELIVERY_TIME_FMT) {
		t.FailNow()
	}

	if wrapper["template_name"] != e.TemplateId {
		t.FailNow()
	}

	ctx := wrapper["template_context"].([]*mandrillTemplateContext)

	if ctx[0].Name != "hello" {
		t.FailNow()
	}

	if ctx[0].Content != "world" {
		t.FailNow()
	}
}

func TestBackendEmail(t *testing.T) {
	e := testutils.TestEmail()

	me, err := b.mandrillEmailForEmail(e)
	if err != nil {
		t.FailNow()
	}

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

func TestBackendAttachments(t *testing.T) {
	e := testutils.TestEmail()

	attachment := testutils.TestAttachment(t)
	e.Attachments = append(e.Attachments, attachment)

	me, err := b.mandrillEmailForEmail(e)
	if err != nil {
		t.FailNow()
	}

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
