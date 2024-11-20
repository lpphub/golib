package render

import (
	"fmt"
	"github.com/pkg/errors"
)

type Error struct {
	Code int
	Msg  string
}

func (err Error) Error() string {
	return err.Msg
}

func (err Error) Wrap(core error) error {
	if core == nil {
		return err
	}
	msg := err.Msg
	err.Msg = core.Error()
	return errors.Wrap(err, msg)
}

func (err Error) WithMessage(msg string) error {
	return errors.WithMessage(err, msg)
}

func (err Error) WithMessagef(format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

func (err Error) Sprintf(v interface{}) Error {
	err.Msg = fmt.Sprintf("%s: %v", err.Msg, v)
	return err
}
