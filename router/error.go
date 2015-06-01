package router

type Error struct{
  Error error
  Conflict bool
  Tail string
  Data map[string]interface{}
}

func NewError (e error, tail ...string) (err *Error) {
  err = &Error{Error: e}
  if len(tail) > 0 {
    err.Tail = tail[0]
  }
  return
}
