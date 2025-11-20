package rfcquery

import "fmt"

type Position struct {
	Offset int
}

type Error struct {
	Pos Position
	Msg string
}

func (e *Error) Error() string {
	return fmt.Sprintf("rfcquery: %s at position %d", e.Msg, e.Pos.Offset)
}

func newError(pos int, format string, args ...any) *Error {
	return &Error{
		Pos: Position{Offset: pos},
		Msg: fmt.Sprintf(format, args...),
	}
}
