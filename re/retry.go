package re

import (
	"context"
	"time"
)

type TryPolicy interface {
	SleepDuration(attempt int, previousSleep time.Duration) time.Duration
	Continue() bool
	WithFactor(factor int64) TryPolicy
	WithInterval(interval time.Duration) TryPolicy
	WithMaxInterval(interval time.Duration) TryPolicy
	WithTimeout(timeout time.Duration) TryPolicy
}

// Try - from app called as `re.Try(ctx, fn)` - retries the given function until it succeeds or time is exceeded.
// TryPolicy specifies how waiting is processed and the timeout.
func Try[T any](ctx context.Context, fn func() (T, error), rp ...TryPolicy) (T, error) {
	var zero T
	var err error
	var result T

	if len(rp) == 0 {
		rp = []TryPolicy{Const()}
	}
	policy := rp[0] // only the first policy is used to make this argument optional
	attempt := 0
	toSleep := time.Duration(0)
	lastSleep := time.Duration(0)

	for policy.Continue() {
		if toSleep > 0 {
			err = Sleep(ctx, toSleep)
			if err != nil {
				return zero, err
			}
			lastSleep = toSleep
		}
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}
		attempt++
		result, err = fn()
		if err == nil {
			return result, nil
		}
		toSleep = policy.SleepDuration(attempt, lastSleep)
	}

	return zero, err
}

type Policy struct {
	Interval    time.Duration
	MaxInterval time.Duration
	Factor      int64
	StopAt      time.Time
	Self        TryPolicy
}

func (b *Policy) WithFactor(factor int64) TryPolicy {
	b.Factor = factor
	return b.Self
}

func (b *Policy) WithInterval(interval time.Duration) TryPolicy {
	b.Interval = interval
	return b.Self
}

func (b *Policy) WithMaxInterval(interval time.Duration) TryPolicy {
	b.MaxInterval = interval
	return b.Self
}

func (b *Policy) WithTimeout(timeout time.Duration) TryPolicy {
	b.StopAt = time.Now().Add(timeout)
	return b.Self
}

func (b *Policy) SleepDuration(_ int, _ time.Duration) time.Duration {
	return b.Interval
}

func (b *Policy) Continue() bool {
	return time.Now().Before(b.StopAt)
}

// Sleep pauses the current goroutine until d has passed or the context is canceled—whichever happens first.
func Sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
