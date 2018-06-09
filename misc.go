package ctx

import (
	"errors"
	"log"
	"time"
)

type lock struct{}

func persistWithHealthCheck(
	maxAttempt int,
	timeout time.Duration,
	m map[string]interface{},
	name string,
	fn func() (interface{}, error),
) interface{} {
	if value, found := m[name]; found {
		return checkForLockOrReturnValue(value)
	}
	m[name] = lock{}

	attempt := 0

	for {
		value, err := fn()
		if nil == err {
			m[name] = value
			return value
		}

		attempt++
		if attempt >= maxAttempt {
			log.Panic(err)
			break
		}
		time.Sleep(timeout)
	}

	return nil
}

func persist(m map[string]interface{}, name string, fn func() interface{}) interface{} {
	if value, found := m[name]; found {
		return checkForLockOrReturnValue(value)
	}
	m[name] = lock{}

	value := fn()
	m[name] = value

	return value
}

const errMsg = "Already set! Use a different name please!"

func panicOnFound(found bool) {
	if found {
		panic(errors.New(errMsg))
	}
}

func checkForLockOrReturnValue(value interface{}) (rtnValue interface{}) {
	switch value := value.(type) {
	case lock:
		panic(errors.New(errMsg))
	default:
		rtnValue = value
	}
	return
}
