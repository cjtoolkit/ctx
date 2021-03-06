package ctx

import (
	"errors"
	"sync"
	"testing"
)

func TestPersistWithHealthCheck(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		m := &sync.Map{}
		name := "test"

		persistWithHealthCheck(2, 0, m, name, func() (interface{}, error) {
			return "valid", nil
		})

		if v, _ := m.Load(name); v.(string) != "valid" {
			t.Error("Should be 'valid'")
		}
	})

	t.Run("Has error on first, none on second", func(t *testing.T) {
		m := &sync.Map{}
		name := "test"

		attempt := -1
		persistWithHealthCheck(2, 0, m, name, func() (interface{}, error) {
			attempt++
			if attempt == 1 {
				return "valid", nil
			}
			return nil, errors.New("I am error")
		})

		if v, _ := m.Load(name); v.(string) != "valid" {
			t.Error("Should be 'valid'")
		}
	})

	t.Run("Reach maxAttempt", func(t *testing.T) {
		m := &sync.Map{}
		name := "test"

		defer func() {
			if recover() == nil {
				t.Error("Recover should be nil.")
			}
		}()

		attempt := -1
		persistWithHealthCheck(2, 0, m, name, func() (interface{}, error) {
			attempt++
			if attempt == 2 {
				return "valid", nil
			}
			return nil, errors.New("I am error")
		})
	})
}

func TestPersist(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		name := "test"
		m := &sync.Map{}

		persist(m, name, func() interface{} {
			return "set"
		})

		if value, _ := m.Load(name); value.(string) != "set" {
			t.Error("Should be 'set'")
		}
	})

	t.Run("Get", func(t *testing.T) {
		name := "test"
		m := &sync.Map{}
		m.Store(name, "get")

		persist(m, name, func() interface{} {
			return "set"
		})

		if value, _ := m.Load(name); value.(string) != "get" {
			t.Error("Should be 'get'")
		}
	})
}

func TestPanicIfFound(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("Recover should be nil.")
			}
		}()
		panicOnFound(true)
	})

	t.Run("Not Found", func(t *testing.T) {
		panicOnFound(false)
	})
}

func TestCheckForLockOrReturnValue(t *testing.T) {
	t.Run("Is Locked", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("Recover should not be nil")
			}
		}()

		checkForLockOrReturnValue(lock{})
	})

	t.Run("Is Unlocked", func(t *testing.T) {
		checkForLockOrReturnValue("hello")
	})
}
