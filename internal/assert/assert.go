package assert

import (
	"fmt"
	"log"
	"os"
)

func Assert(condition bool, message any) {
	if condition {
		return
	}

	switch msg := message.(type) {
	case error:
		log.Println(msg.Error())
	default:
		log.Println(msg)
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
