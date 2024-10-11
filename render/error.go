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

func (err Error) Cause() error {
	return errors.Cause(err)
}

func (err Error) Wrap(core error) error {
	if core == nil {
		return err
	}
	msg := err.Msg
	err.Msg = core.Error()
	return errors.Wrap(err, msg)
}

func (err Error) WithMessage(msg string) Error {
	err.Msg = fmt.Sprintf("%s: %s", err.Msg, msg)
	return err
}

func (err Error) Sprintf(format string, args ...interface{}) Error {
	return err.WithMessage(fmt.Sprintf(format, args...))
}

func (err Error) Sprintf2(v ...interface{}) Error {
	err.Msg = fmt.Sprintf(err.Msg, v...)
	return err
}
