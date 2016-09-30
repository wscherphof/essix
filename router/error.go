package router

type tailType struct {
	dir  string
	name string
}

// Error is what an ErrorHandler should return instead of executing a template
type Error struct {
	Error    error
	Conflict bool
	Tail     *tailType
	Data     map[string]interface{}
}

// NewError constructs an Error
func NewError(e error, tail ...string) (err *Error) {
	if e != nil {
		err = &Error{Error: e}
		if len(tail) == 2 {
			err.Tail = &tailType{
				dir:  tail[0],
				name: tail[1] + "_error-tail",
			}
		}
	}
	return
}

// IfError returns a new Error constructed from e, or nil if e was nil
func IfError(e error, tail ...string) (err *Error) {
	return NewError(e, tail...)
}
