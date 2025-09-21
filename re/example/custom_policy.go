//go:build tests

package example

import (
	"time"

	"github.com/digitalstraw/retry-go/re"
)

type CustomPolicy struct {
	*re.TryBasePolicy

	failures int
}

var _ re.TryPolicy = (*CustomPolicy)(nil)

func Custom() *CustomPolicy {
	c := &CustomPolicy{}
	c.TryBasePolicy = &re.TryBasePolicy{
		Interval:    1 * time.Second,
		MaxInterval: 7 * time.Second,
		StopAt:      re.StopAt(10 * time.Second), //nolint:mnd // Unimportant
		Self:        c,
	}
	return c
}

// SleepDuration produces following durations in default: 3s, 5s, 7s, timeout.
func (c *CustomPolicy) SleepDuration(attempt int, _ time.Duration) time.Duration {
	c.failures++
	return c.Interval + time.Duration(attempt)*2*time.Second // Custom logic
}

func (c *CustomPolicy) Failures() int {
	return c.failures
}
