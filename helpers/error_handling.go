package helpers

import (
	"github.com/go-errors/errors"
	"os"
)

func PanicOnError(err error) {
	if err == nil {
		return
	}

	if e, ok := err.(*errors.Error); ok {
		_, er := os.Stderr.WriteString(e.ErrorStack())
		if er != nil {
			panic(er)
		}
	}
	panic(err)
}
