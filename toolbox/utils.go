package toolbox

import (
	"errors"
	"log/slog"
)

func Uups(err error) {
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func MustBeTrue(cond bool, msg string, args ...interface{}) {
	if !cond {
		slog.Error("msg", args...)
		panic(errors.New(msg))
	}
}
