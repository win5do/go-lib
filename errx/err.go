package errx

import (
	"sync/atomic"

	errors2 "github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors2.StackTrace
}

var stackFlag uint32 = 1

func SetFlag(b bool) {
	var x uint32 = 1
	if !b {
		x = 0
	}

	atomic.StoreUint32(&stackFlag, x)
}

func WithStackOnce(err error) error {
	if stackFlag == 0 {
		return err
	}

	_, ok := err.(stackTracer)
	if ok {
		return err
	}

	return errors2.WithStack(err)
}
