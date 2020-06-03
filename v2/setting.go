package ctx

import (
	"time"
)

type Settings struct {
	// Set Health Check Max Attempt on Background Context
	MaxAttempt int

	// Set Health Check Time Out on Background Context
	Timeout time.Duration
}

type settingsKey struct{}

func (c *contextBase) settings() Settings {
	v, found := c.get(settingsKey{})
	if found {
		v := v.(*Settings)
		return *v
	}
	s := defaultSettings()
	c.set(settingsKey{}, s)
	return *s
}

func defaultSettings() *Settings {
	return &Settings{
		MaxAttempt: 1,
		Timeout:    0,
	}
}

/*
Get Settings from Context
*/
func ContextSettings(context Context) *Settings {
	v, found := context.Get(settingsKey{})
	if found {
		return v.(*Settings)
	}
	s := defaultSettings()
	context.Set(settingsKey{}, s)
	return s
}
