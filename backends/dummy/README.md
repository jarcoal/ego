#Dummy

The dummy backend is designed for development/testing; it doesn't send any emails, it simply logs the contents of an email to a given logger.

#Example

```go
package main

import (
	"net/mail"
	"github.com/jarcoal/ego"
	"github.com/jarcoal/ego/dummy"
	"log"
)

func main() {
	backend := dummy.NewBackend(log.Printf)

	email := ego.NewEmail()
	email.From = &mail.Address{"Jane Smith", "jane@smith.com"}
	email.AddRecipient("John Smith", "john@smith.com")
	email.Subject = "Hello World"
	email.HTMLBody = "<h1>Hello World</h1>"

	backend.SendEmail(email)
}
```