#SendGrid

Backend for the promotional and transactional email provider [SendGrid](http://sendgrid.com/).

#Support

SendGrid supports most features, despite their incredibly crappy API, docs, and admin panel.

* Templating (very limited)
* Tagging
* Attachments
* Click/Open tracking

#Example

```go
package main

import (
	"net/mail"
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/sendgrid"
)

func main() {
	backend := sendgrid.NewBackend("<my-sendgrid-username>", "<my-sendgrid-password>")

	email := ego.NewEmail()
	email.From = &mail.Address{"Jane Smith", "jane@smith.com"}
	email.AddRecipient("John Smith", "john@smith.com")
	email.Subject = "Hello World"
	email.HTMLBody = "<h1>Hello World</h1>"

	backend.DispatchEmail(email)
}
```