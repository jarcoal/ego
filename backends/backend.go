package backends

import (
	"github.com/jarcoal/ego"
)

// Backend is the interface that all backends must implement to send emails
type Backend interface {
	SendEmail(*ego.Email) error
}
