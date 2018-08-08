package ctx

import (
	"errors"
	"log"
	"sync"
	"time"
)

type lock struct{}

func persistWithHealthCheck(
	maxAttempt int,
	timeout time.Duration,
	m *sync.Map,
	key interface{},
	fn func() (interface{}, error),
) interface{} {
	if value, found := m.Load(key); found {
		return checkForLockOrReturnValue(value)
	}
	m.Store(key, lock{})

	attempt := 0

	for {
		value, err := fn()
		if nil == err {
			m.Store(key, value)
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

func persist(m *sync.Map, key interface{}, fn func() interface{}) interface{} {
	if value, found := m.Load(key); found {
		return checkForLockOrReturnValue(value)
	}
	m.Store(key, lock{})

	value := fn()
	m.Store(key, value)

	return value
}

const errMsg = "Collision Detected!"

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
