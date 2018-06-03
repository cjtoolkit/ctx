package ctx

import (
	"errors"
	"log"
	"time"
)

func persistWithHealthCheck(
	maxAttempt int,
	timeout time.Duration,
	m map[string]interface{},
	name string,
	fn func() (interface{}, error),
) interface{} {
	if value, found := m[name]; found {
		return value
	}

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
		return value
	}

	value := fn()
	m[name] = value

	return value
}

func panicOnFound(found bool) {
	if found {
		panic(errors.New("Already set! Use a different name please!"))
	}
}
