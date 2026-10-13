package re

import "errors"

// permanentError marks an error a retry cannot fix.
type permanentError struct {
	err error
}

func (p *permanentError) Error() string {
	return p.err.Error()
}

func (p *permanentError) Unwrap() error {
	return p.err
}

// Permanent marks err as one a retry cannot fix: when fn returns it, Try stops
// at once and returns err itself, unwrapped. Permanent(nil) is nil.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &permanentError{err: err}
}

// asPermanent returns the error wrapped by Permanent, if err carries one.
func asPermanent(err error) (error, bool) {
	var p *permanentError
	if !errors.As(err, &p) {
		return nil, false
	}
	return p.err, true
}
