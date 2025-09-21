package user

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
		MaxInterval: 5 * time.Second,
		StopAt:      re.StopAt(10 * time.Second), //nolint:mnd // Unimportant
		Self:        c,
	}
	return c
}

func (c *CustomPolicy) SleepDuration(attempt int, _ time.Duration) time.Duration {
	c.failures++
	return c.Interval + time.Duration(attempt)*2*time.Second // Custom logic: increase sleep by 100ms each attempt
}

func (c *CustomPolicy) Failures() int {
	return c.failures
}
