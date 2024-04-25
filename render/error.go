package render

import "fmt"

type Error struct {
	Code int
	Msg  string
}

func (err Error) Error() string {
	return err.Msg
}

func (err Error) Case() string {
	return err.Msg
}

func (err Error) Sprintf(v ...interface{}) Error {
	err.Msg = fmt.Sprintf(err.Msg, v...)
	return err
}
