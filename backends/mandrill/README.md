#Mandrill

Backend for the transactional email provider [Mandrill](http://mandrill.com/).

#Support

Mandrill is a pretty full-featured sending service, so there isn't much they don't support.

* Templating
* Delayed Delivery
* Tagging
* Attachments
* Click/Open tracking

#Example

```go
package main

import (
	"net/mail"
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/mandrill"
)

func main() {
	backend := mandrill.NewBackend("<my-mandrill-api-key>")

	email := ego.NewEmail()
	email.From = &mail.Address{"Jane Smith", "jane@smith.com"}
	email.AddRecipient("John Smith", "john@smith.com")
	email.Subject = "Hello World"
	email.HTMLBody = "<h1>Hello World</h1>"

	backend.SendEmail(email)
}
```