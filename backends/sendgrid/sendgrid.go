// SendGrid transactional email backend
//
// Website: http://sendgrid.com/
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

const apiURL = "https://sendgrid.com/api/mail.send.json"

var _ backends.Backend = (*sendGridBackend)(nil)

// NewBackend creates a new SendGrid backend that is bound to the given credentials.
func NewBackend(username, password string) backends.Backend {
	return &sendGridBackend{username, password}
}

type sendGridBackend struct {
	username, password string
}

func (s *sendGridBackend) SendEmail(e *ego.Email) error {
	// get the parameters we're going to be posting to sendgrid
	params, err := s.paramsForEmail(e)
	if err != nil {
		return err
	}

	// perform the request
	resp, err := http.PostForm(apiURL, params)
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
	params.Set("html", e.HTMLBody)
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

	// add any headers
	if len(e.Headers) > 0 {
		headerMap := make(map[string]string)

		for header := range e.Headers {
			headerMap[header] = e.Headers.Get(header)
		}

		headerMapEncoded, err := json.Marshal(headerMap)
		if err != nil {
			return nil, fmt.Errorf("failed to encode headers: %s", err)
		}

		params.Set("headers", string(headerMapEncoded))
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
	xSMTPApiParams := make(map[string]interface{})

	// apply the tags (which sendgrid calls categories)
	if len(e.Tags) > 0 {
		xSMTPApiParams["category"] = e.Tags
	}

	// template filter properties
	if e.TemplateID != "" {
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
		xSMTPApiParams["filters"] = map[string]interface{}{
			"templates": map[string]interface{}{
				"settings": map[string]interface{}{
					"enabled":     1,
					"template_id": e.TemplateID,
				},
			},
		}
	}

	// x-smtp params are to be json encoded before adding to the params
	xSMTPAPIEncoded, err := json.Marshal(xSMTPApiParams)
	if err != nil {
		return nil, errors.New("failed to encode x-smtpapi parameters")
	}
	params.Set("x-smtpapi", string(xSMTPAPIEncoded))

	return params, nil
}
