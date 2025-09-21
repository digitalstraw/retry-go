package re

import "time"

type ConstPolicy struct {
	*TryBasePolicy
}

var _ TryPolicy = (*ConstPolicy)(nil)

// Const represents constant policy. It delivers a single Interval of stable length.
// In default: waits for 1 second between each attempt and 30 seconds timeout.
func Const() *ConstPolicy {
	interval := 1 * time.Second
	cp := &ConstPolicy{}
	cp.TryBasePolicy = &TryBasePolicy{
		Interval:    interval,
		MaxInterval: interval,
		StopAt:      StopAt(defaultConstTimeout),
		Self:        cp,
	}
	return cp
}
