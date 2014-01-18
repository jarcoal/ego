// PostageApp
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

const POSTAGEAPP_API_URL = "https://api.postageapp.com/v.1.0/send_message.json"

var _ backends.Backend = (*postageAppBackend)(nil)

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

	resp, err := http.Post(POSTAGEAPP_API_URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to post to postageapp: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		postageAppErr := &postageAppError{}

		if err := json.NewDecoder(resp.Body).Decode(postageAppErr); err != nil {
			return fmt.Errorf("received %s from postageapp and couldn't decode error payload: %s",
				resp.StatusCode, err)
		}

		return postageAppErr
	}

	return nil
}

func (p *postageAppBackend) wrapperForEmail(e *ego.Email) (*postageAppWrapper, error) {
	pa := &postageAppArguments{
		Recipients: make([]string, 0),
		Headers: map[string]string{
			"subject": e.Subject,
			"from":    e.From.String(),
		},
		Content: map[string]string{
			"text/plain": e.TextBody,
			"text/html":  e.HtmlBody,
		},
	}

	// recipients
	for _, recip := range e.To {
		pa.Recipients = append(pa.Recipients, recip.String())
	}

	// reply to
	if e.ReplyTo != nil {
		pa.Headers["reply-to"] = e.ReplyTo.String()
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
	if e.TemplateId != "" {
		pa.Template = e.TemplateId
		pa.Variables = e.TemplateContext
	}

	// put it in the wrapper
	wrapper := &postageAppWrapper{
		ApiKey:    p.apiKey,
		Arguments: pa,
	}

	return wrapper, nil
}

type postageAppWrapper struct {
	ApiKey    string               `json:"api_key"`
	Uid       string               `json:"uid,omitempty"`
	Arguments *postageAppArguments `json:"arguments"`
}

type postageAppArguments struct {
	Recipients        []string                         `json:"recipients"`
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
	Uid     string `json:"uid"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (p *postageAppError) Error() string {
	return fmt.Sprintf("%v %v - %v", p.Uid, p.Status, p.Message)
}