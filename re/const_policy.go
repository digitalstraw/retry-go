package re

import "time"

const (
	defaultConstTimeout = 30 * time.Second
)

type ConstPolicy struct {
	*Policy
}

var _ TryPolicy = (*ConstPolicy)(nil)

// Const represents constant policy. It delivers a single Interval of stable length.
// In default: waits for 1 second between each attempt and 30 seconds timeout.
func Const() *ConstPolicy {
	interval := 1 * time.Second
	cp := &ConstPolicy{}
	cp.Policy = &Policy{
		Interval: interval,
		StopAt:   time.Now().Add(defaultConstTimeout),
		Self:     cp,
	}
	return cp
}
