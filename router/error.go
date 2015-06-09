package router

type tailType struct {
	dir  string
	name string
}

type Error struct {
	Error    error
	Conflict bool
	Tail     *tailType
	Data     map[string]interface{}
}

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

func IfError(e error, tail ...string) (err *Error) {
	return NewError(e, tail...)
}
