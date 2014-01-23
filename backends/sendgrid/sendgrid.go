// SendGrid
// API Docs: http://sendgrid.com/docs/API_Reference/Web_API/mail.html
package sendgrid

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/backends"
	"io/ioutil"
	"net/http"
	"net/url"
)

const SENDGRID_API_URL = "https://sendgrid.com/api/mail.send.json"

var _ backends.Backend = (*sendGridBackend)(nil)

func NewBackend(username, password string) backends.Backend {
	return &sendGridBackend{username, password}
}

type sendGridBackend struct {
	username, password string
}

func (s *sendGridBackend) DispatchEmail(e *ego.Email) error {
	// get the parameters we're going to be posting to sendgrid
	params, err := s.paramsForEmail(e)
	if err != nil {
		return err
	}

	// perform the request
	resp, err := http.PostForm(SENDGRID_API_URL, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("received bad status code from sendgrid: " + resp.Status)
	}

	return nil
}

func (s *sendGridBackend) paramsForEmail(e *ego.Email) (url.Values, error) {
	params := url.Values{}

	// apply our credentials
	params.Set("api_user", s.username)
	params.Set("api_key", s.password)

	// general information
	params.Set("subject", e.Subject)
	params.Set("text", e.TextBody)
	params.Set("html", e.HtmlBody)
	params.Set("from", e.From.Address)
	params.Set("fromname", e.From.Name)

	if e.ReplyTo != nil {
		params.Set("replyto", e.ReplyTo.Address)
	}

	for _, to := range e.To {
		params.Add("to[]", to.Email.Address)
		params.Add("toname[]", to.Email.Name)
	}

	for _, bcc := range e.Bcc {
		params.Add("bcc[]", bcc.Email.Address)
	}

	// add any attachments
	for _, attachment := range e.Attachments {
		attachmentBytes, err := ioutil.ReadAll(attachment.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to read attachment %s", attachment.Name)
		}
		params.Set(fmt.Sprintf("files[%v]", attachment.Name), string(attachmentBytes))
	}

	// these are misc parameters that get fed into the smtp api
	xSmtpApiParams := make(map[string]interface{})

	// apply the tags (which sendgrid calls categories)
	if len(e.Tags) > 0 {
		xSmtpApiParams["category"] = e.Tags
	}

	// template filter properties
	if e.TemplateId != "" {
		// there should be only one template context variable named 'body'
		if len(e.TemplateContext) != 1 {
			return nil, errors.New("template context for SendGrid can only contain one variable 'body'")
		}
		bodyCtx, ok := e.TemplateContext["body"]
		if !ok {
			return nil, errors.New("template context for SendGrid can only contain one variable 'body'")
		}

		params.Set("body", bodyCtx)

		// apply the template id
		xSmtpApiParams["filters"] = map[string]interface{}{
			"templates": map[string]interface{}{
				"settings": map[string]interface{}{
					"enabled":     1,
					"template_id": e.TemplateId,
				},
			},
		}
	}

	// x-smtp params are to be json encoded before adding to the params
	xSmtpApiEncoded, err := json.Marshal(xSmtpApiParams)
	if err != nil {
		return nil, errors.New("failed to encode x-smtpapi parameters")
	}
	params.Set("x-smtpapi", string(xSmtpApiEncoded))

	return params, nil
}
