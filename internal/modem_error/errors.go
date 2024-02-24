package modem_error

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ConnNotFound Error = "connection not found"
	ConnAlreadyInUse Error = "connection already in use"
	ConnIdOutOfRange Error = "connection id out of range"
	RouteNotFound Error = "route not found"
	RouterAlreadyExists Error = "route already exists"
)
