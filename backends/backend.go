package backends

import (
	"github.com/jarcoal/ego"
)

// Backend is the interface that all backends must implement to send emails
type Backend interface {
	DispatchEmail(*ego.Email) error
}
