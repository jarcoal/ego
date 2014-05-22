#PostageApp

Backend for the transactional email provider [PostageApp](http://postageapp.com/).

#Support

* Templating
* Tagging
* Attachments
* Click/Open tracking

#Example

```go
package main

import (
	"net/mail"
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/postageapp"
)

func main() {
	backend := postageapp.NewBackend("<my-postageapp-api-key>")

	email := ego.NewEmail()
	email.From = &mail.Address{"Jane Smith", "jane@smith.com"}
	email.AddRecipient("John Smith", "john@smith.com")
	email.Subject = "Hello World"
	email.HTMLBody = "<h1>Hello World</h1>"

	backend.SendEmail(email)
}
```