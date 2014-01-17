package ego

type Backend interface {
	DispatchEmail(*Email) error
}
