package myfdbstorage

// Error type with internal
type Error struct {
	Msg           string
	InternalError error
}

func (e *Error) Error() string {
	if e.InternalError == nil {
		return e.Msg
	}
	return e.Msg + "\n" + e.InternalError.Error()
}

// ErrorNew - Create new error
func ErrorNew(msg string, internalError error) error {
	return &Error{Msg: msg, InternalError: internalError}
}
