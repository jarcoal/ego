// Mandrill
// API Docs: https://mandrillapp.com/api/docs/messages.JSON.html
package backends

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/backends"
	"io/ioutil"
	"net/http"
)

const MANDRILL_DELIVERY_TIME_FMT = "2006-01-02T15:04:05"
const MANDRILL_API_URL = "https://mandrillapp.com/api/1.0/messages/%s.json"

var _ backends.Backend = (*MandrillBackend)(nil)

func NewMandrillBackend(apiKey string) *MandrillBackend {
	return &MandrillBackend{apiKey}
}

type MandrillBackend struct {
	apiKey string
}

func (m *MandrillBackend) DispatchEmail(e *ego.Email) error {
	// convert the email to a mandrillEmail struct that will be json-serialized and sent out
	me, err := m.mandrillEmailForEmail(e)
	if err != nil {
		return fmt.Errorf("failed to build mandrill email: %s", err)
	}

	// wrap the mandrill email and encode it
	body, err := json.Marshal(m.wrapMandrillEmail(e, me))
	if err != nil {
		return fmt.Errorf("failed to encode mandrill payload: %s", err)
	}

	// mandrill uses different endpoints if you're sending a templated email
	var apiUrl string
	if e.TemplateId != "" {
		apiUrl = fmt.Sprintf(MANDRILL_API_URL, "send-template")
	} else {
		apiUrl = fmt.Sprintf(MANDRILL_API_URL, "send")
	}

	// make the request to mandrill's api
	resp, err := http.Post(apiUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to post to mandrill: %s", err)
	}
	defer resp.Body.Close()

	// if we got a bad status code, read out the error body
	if resp.StatusCode != 200 {
		mandrillErr := &mandrillError{}

		if err := json.NewDecoder(resp.Body).Decode(err); err != nil {
			return fmt.Errorf("received %s from mandrill and couldn't decode error payload: %s",
				resp.StatusCode, err)
		}

		return mandrillErr
	}

	return nil
}

func (m *MandrillBackend) mandrillEmailForEmail(e *ego.Email) (*mandrillEmail, error) {
	me := &mandrillEmail{
		To:                 make([]*mandrillRecipient, 0, len(e.To)),
		Attachments:        make([]*mandrillAttachment, 0, len(e.Attachments)),
		Headers:            make(map[string]string),
		Html:               e.HtmlBody,
		Text:               e.TextBody,
		Subject:            e.Subject,
		FromEmail:          e.From.Address,
		FromName:           e.From.Name,
		TrackOpens:         e.TrackOpens,
		TrackClicks:        e.TrackClicks,
		Tags:               e.Tags,
		Subaccount:         e.SubAccount,
		PreserveRecipients: e.VisibleRecipients,
	}

	if e.ReplyTo != nil {
		me.Headers["Reply-To"] = e.ReplyTo.Address
	}

	// assign the recipients
	for _, to := range e.To {
		me.To = append(me.To, &mandrillRecipient{"to", to.Name, to.Address})
	}

	// add attachments
	for _, attachment := range e.Attachments {
		attachmentBytes, err := ioutil.ReadAll(attachment.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s attachment: %s", attachment.Name, err)
		}

		me.Attachments = append(me.Attachments, &mandrillAttachment{
			attachment.Mimetype,
			attachment.Name,
			base64.StdEncoding.EncodeToString(attachmentBytes),
		})
	}

	return me, nil
}

func (m *MandrillBackend) wrapMandrillEmail(e *ego.Email, me *mandrillEmail) map[string]interface{} {
	wrapper := map[string]interface{}{
		"key":     m.apiKey,
		"message": me,
	}

	// if a service-side template has been specified
	if e.TemplateId != "" {
		wrapper["template_name"] = e.TemplateId
		wrapper["template_context"] = mandrillCtxFromContext(e.TemplateContext)
	}

	// assign scheduled send date if provided
	if !e.DeliveryTime.IsZero() {
		wrapper["send_at"] = e.DeliveryTime.Format(MANDRILL_DELIVERY_TIME_FMT)
	}

	return wrapper
}

// mandrillCtxFromContext converts the context map into a slice of mandrillTemplateContext
func mandrillCtxFromContext(context map[string]string) []*mandrillTemplateContext {
	resp := make([]*mandrillTemplateContext, 0, len(context))

	for k, v := range context {
		resp = append(resp, &mandrillTemplateContext{
			Name:    k,
			Content: v,
		})
	}
	return resp
}

// mandrillEmail is an email message that can be serialized and delivered to mandrill's web api
type mandrillEmail struct {
	To                 []*mandrillRecipient  `json:"to"`
	Attachments        []*mandrillAttachment `json:"attachments,omitempty"`
	Html               string                `json:"html,omitempty"`
	Text               string                `json:"text,omitempty"`
	Subject            string                `json:"subject"`
	FromEmail          string                `json:"from_email"`
	FromName           string                `json:"from_name"`
	TrackOpens         bool                  `json:"track_opens"`
	TrackClicks        bool                  `json:"track_clicks"`
	Tags               []string              `json:"tags,omitempty"`
	Subaccount         string                `json:"subaccount,omitempty"`
	PreserveRecipients bool                  `json:"preserve_recipients"`
	Headers            map[string]string     `json:"headers,omitempty"`
}

// mandrillTemplateContext represents a single key/value pair to be used in
// mandrill's templating system
type mandrillTemplateContext struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// mandrillRecipient represents a single recipient in a mandrill email
type mandrillRecipient struct {
	Type  string `json:"type"`  // always set to 'to'
	Name  string `json:"name"`  // recipient's name
	Email string `json:"email"` // recipient's address
}

// mandrillAttachment represents a single attachment in a mandrill email
type mandrillAttachment struct {
	Type    string `json:"type"`    // mimetype of the attachment
	Name    string `json:"name"`    // file name of the attachment
	Content string `json:"content"` // base64-encoded version of the file
}

// mandrillError represents a json-encoded error returned by mandrill from an api call
type mandrillError struct {
	Status  string `json:status`
	Code    int    `json:code`
	Name    string `json:name`
	Message string `json:message`
}

func (m *mandrillError) Error() string {
	return fmt.Sprintf("%v %v - %v", m.Code, m.Name, m.Message)
}