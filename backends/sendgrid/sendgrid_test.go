package backends

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/ego/testutils"
	"io"
	"io/ioutil"
	"net/url"
	"testing"
)

var b = sendGridBackend{"test-username", "test-password"}

func TestSendGridBackendGeneral(t *testing.T) {
	e := testutils.TestEmail()

	params, err := b.paramsForEmail(e)
	if err != nil {
		t.FailNow()
	}

	if params.Get("subject") != e.Subject {
		t.FailNow()
	}

	if params.Get("text") != e.TextBody {
		t.FailNow()
	}

	if params.Get("html") != e.HtmlBody {
		t.FailNow()
	}

	if params.Get("from") != e.From.Address {
		t.FailNow()
	}

	if params.Get("fromname") != e.From.Name {
		t.FailNow()
	}

	if params.Get("replyto") != e.ReplyTo.Address {
		t.FailNow()
	}

	if len(params["to[]"]) != len(e.To) || len(params["toname[]"]) != len(e.To) {
		t.FailNow()
	}
}

func TestSendGridBackendCategories(t *testing.T) {
	e := testutils.TestEmail()

	params, err := b.paramsForEmail(e)
	if err != nil {
		t.FailNow()
	}

	xSmtpApi := decodeXsmtpApi(t, params)
	if len(xSmtpApi["category"].([]interface{})) != len(e.Tags) {
		t.FailNow()
	}
}

func TestSendGridBackendTemplating(t *testing.T) {
	e := testutils.TestEmail()
	e.TemplateContext = make(map[string]string)

	// we haven't added any template info yet
	func() {
		params, err := b.paramsForEmail(e)
		if err != nil {
			t.FailNow()
		}

		// should be no 'filters' in the xsmtpapi params
		xSmtpApi := decodeXsmtpApi(t, params)
		if _, ok := xSmtpApi["filters"]; ok {
			t.FailNow()
		}

		// should be no 'body' param
		if params.Get("body") != "" {
			t.FailNow()
		}
	}()

	// add some invalid template data that should trigger an error
	func() {
		e.TemplateId = "test-template"
		e.TemplateContext["not-body"] = "some invalid variable"

		if _, err := b.paramsForEmail(e); err == nil {
			t.FailNow()
		}
	}()

	// finally some good data
	func() {
		e.TemplateId = "test-template"
		e.TemplateContext = map[string]string{"body": "legit"}

		params, err := b.paramsForEmail(e)
		if err != nil {
			t.FailNow()
		}

		xSmtpApi := decodeXsmtpApi(t, params)
		filters := xSmtpApi["filters"].(map[string]interface{})
		templates := filters["templates"].(map[string]interface{})
		settings := templates["settings"].(map[string]interface{})

		if settings["enabled"] != float64(1) {
			t.FailNow()
		}

		if settings["template_id"] != e.TemplateId {
			t.FailNow()
		}
	}()
}

func TestSendGridBackendAttachments(t *testing.T) {
	e := testutils.TestEmail()

	attachment := testutils.TestAttachment(t)
	e.Attachments = append(e.Attachments, attachment)

	params, err := b.paramsForEmail(e)
	if err != nil {
		t.FailNow()
	}

	// the sendgrid backend will have read the contents of the reader,
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
	if params.Get(fmt.Sprintf("files[%v]", e.Attachments[0].Name)) != string(fileData) {
		t.FailNow()
	}
}

func decodeXsmtpApi(t *testing.T, params url.Values) map[string]interface{} {
	xSmtpApi := make(map[string]interface{})

	if err := json.Unmarshal([]byte(params.Get("x-smtpapi")), &xSmtpApi); err != nil {
		t.FailNow()
	}

	return xSmtpApi
}
