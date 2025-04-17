package assert

import (
	"fmt"
	"os"
)

func Assert(condition bool, message any) {
	if condition {
		return
	}

	switch msg := message.(type) {
	case error:
		fmt.Fprintln(os.Stderr, msg.Error())
	default:
		fmt.Fprintln(os.Stderr, msg)
	}

	os.Exit(1)
}

func Assertf(condition bool, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	Assert(condition, message)
}

func AssertE(err error) {
	Assert(err == nil, err)
}
