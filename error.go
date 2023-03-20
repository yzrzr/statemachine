package statemachine

import (
	"errors"
)

func NewError(msg string) error {
	return Error{
		msg: msg,
	}
}

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}

func IsStateMachineError(err error) bool {
	e := &Error{}
	return errors.As(err, e)
}
