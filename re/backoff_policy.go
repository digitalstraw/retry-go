package re

import "time"

const (
	defaultBackoffTimeout     = 2 * time.Minute
	defaultBackoffInterval    = 100 * time.Millisecond
	defaultBackoffMaxInterval = 60 * time.Second
	defaultBackoffFactor      = 2
)

type BackoffPolicy struct {
	*Policy
}

var _ TryPolicy = (*BackoffPolicy)(nil)

// Backoff increases waiting time between attempts by multiplying previous delay time with factor.
// It is limited by MaxInterval which represents maximum time between two attempts.
// Also limited with timeout.
// In default: factor = 2, 100ms first delay, then 200, 400, 600 etc. up to 60 seconds and then exactly 60 seconds with 2 minutes timeout.
func Backoff() *BackoffPolicy {
	bp := &BackoffPolicy{}
	bp.Policy = &Policy{
		Interval:    defaultBackoffInterval,
		MaxInterval: defaultBackoffMaxInterval,
		Factor:      defaultBackoffFactor,
		StopAt:      time.Now().Add(defaultBackoffTimeout),
		Self:        bp,
	}
	return bp
}

func (b *BackoffPolicy) SleepDuration(attempt int, previousSleep time.Duration) time.Duration {
	if attempt == 1 {
		return b.Interval
	}
	sleep := time.Duration(previousSleep.Nanoseconds() * b.Factor)
	if sleep > b.MaxInterval {
		return b.MaxInterval
	}
	return sleep
}
