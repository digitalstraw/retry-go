//go:build tests

package example

import (
	"time"

	"github.com/digitalstraw/retry-go/re"
)

const (
	defaultTimeout = 10 * time.Second
)

type CustomPolicy struct {
	*re.Policy

	failures int
}

var _ re.TryPolicy = (*CustomPolicy)(nil)

func Custom() *CustomPolicy {
	c := &CustomPolicy{}
	c.Policy = &re.Policy{
		Interval: 1 * time.Second,
		StopAt:   time.Now().Add(defaultTimeout),
		Self:     c,
	}
	return c
}

// SleepDuration produces following durations in default: 3s, 5s, 7s, timeout.
func (c *CustomPolicy) SleepDuration(attempt int, _ time.Duration) time.Duration {
	c.failures++
	return c.Interval + time.Duration(attempt*2)*time.Second // Custom logic
}

func (c *CustomPolicy) Failures() int {
	return c.failures
}
