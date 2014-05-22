// PostageApp
//
// Website: http://postageapp.com/
// API Docs: http://help.postageapp.com/kb/api/send_message

package postageapp

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

const apiURL = "https://api.postageapp.com/v.1.0/send_message.json"

var _ backends.Backend = (*postageAppBackend)(nil)

// NewBackend returns a Postageapp backend bound to the API key
func NewBackend(apiKey string) backends.Backend {
	return &postageAppBackend{apiKey}
}

type postageAppBackend struct {
	apiKey string
}

func (p *postageAppBackend) DispatchEmail(e *ego.Email) error {
	wrapper, err := p.wrapperForEmail(e)
	if err != nil {
		return fmt.Errorf("failed to build postageapp wrapper: %s", err)
	}

	body, err := json.Marshal(wrapper)
	if err != nil {
		return fmt.Errorf("failed to encode postageapp payload: %s", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to post to postageapp: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		postageAppErr := &postageAppError{}

		if err := json.NewDecoder(resp.Body).Decode(postageAppErr); err != nil {
			return fmt.Errorf("received %s from postageapp and couldn't decode error payload: %s",
				resp.Status, err)
		}

		return postageAppErr
	}

	return nil
}

func (p *postageAppBackend) wrapperForEmail(e *ego.Email) (*postageAppWrapper, error) {
	pa := &postageAppArguments{
		Recipients: make(map[string]map[string]string),
		Headers: map[string]string{
			"subject": e.Subject,
			"from":    e.From.String(),
		},
		Content: map[string]string{
			"text/plain": e.TextBody,
			"text/html":  e.HTMLBody,
		},
	}

	// recipients
	for _, to := range e.To {
		pa.Recipients[to.Email.String()] = to.TemplateContext
	}

	// reply to
	if e.ReplyTo != nil {
		pa.Headers["reply-to"] = e.ReplyTo.String()
	}

	// headers
	if len(e.Headers) > 0 {
		for header := range e.Headers {
			pa.Headers[header] = e.Headers.Get(header)
		}
	}

	// attachments
	if len(e.Attachments) > 0 {
		pa.Attachments = make(map[string]*postageAppAttachment)

		for _, attachment := range e.Attachments {
			data, err := ioutil.ReadAll(attachment.Data)
			if err != nil {
				return nil, fmt.Errorf("failed to read %s attachment: %s", attachment.Name, err)
			}

			pa.Attachments[attachment.Name] = &postageAppAttachment{
				ContentType: attachment.Mimetype,
				Content:     base64.StdEncoding.EncodeToString(data),
			}
		}
	}

	// templating
	if e.TemplateID != "" {
		pa.Template = e.TemplateID
		pa.Variables = e.TemplateContext
	}

	// put it in the wrapper
	wrapper := &postageAppWrapper{
		APIKey:    p.apiKey,
		Arguments: pa,
	}

	return wrapper, nil
}

type postageAppWrapper struct {
	APIKey    string               `json:"api_key"`
	UID       string               `json:"uid,omitempty"`
	Arguments *postageAppArguments `json:"arguments"`
}

type postageAppArguments struct {
	Recipients        map[string]map[string]string     `json:"recipients"`
	Headers           map[string]string                `json:"headers"`
	Content           map[string]string                `json:"content"`
	Attachments       map[string]*postageAppAttachment `json:"attachments"`
	Template          string                           `json:"template"`
	Variables         map[string]string                `json:"variables"`
	RecipientOverride string                           `json:"recipient_override,omitempty"`
}

type postageAppAttachment struct {
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

// postageAppError is a representation of an error response from PostageApp's API
type postageAppError struct {
	UID     string `json:"uid"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (p *postageAppError) Error() string {
	return fmt.Sprintf("%v %v - %v", p.UID, p.Status, p.Message)
}
