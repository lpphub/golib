package render

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestError(t *testing.T) {
	err := &Error{
		Code: 500,
		Msg:  "Internal Server Error",
	}

	e1 := errors.New("aaa")

	e2 := err.Wrap(e1)

	e3 := err.WithMessagef("消息错误： %s", "bbb")

	e4 := err.Sprintf(err).Sprintf(e2)

	fmt.Println(e1.Error())
	fmt.Println(e2.Error())
	fmt.Println(e3.Error())
	fmt.Println(e4.Error())

	fmt.Println("------------------------")

	code, msg := -1, err.Error()
	var err2 Error
	if errors.As(e2, &err2) {
		code = err2.Code
		msg = err2.Error()
	}

	fmt.Printf("%d - %s", code, msg)
}
