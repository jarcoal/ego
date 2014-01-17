package backends

import (
	"github.com/jarcoal/ego"
)

type Backend interface {
	DispatchEmail(*ego.Email) error
}
