package helpers

import (
	"errors"
	"fmt"
	"github.com/ansel1/merry"
	go_errors "github.com/go-errors/errors"
)

func PanicOnError(err error) {
	if err == nil {
		return
	}
	panic(errors.New(ErrorToString(err)))
}

func ErrorToString(err error) string {
	res := fmt.Sprintf("Error '%T'. ", err)
	if e, ok := err.(*go_errors.Error); ok {
		res += e.ErrorStack()
	}
	res += merry.Details(err)
	return res
}
