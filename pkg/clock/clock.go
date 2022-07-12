package clock

import "time"

type Clock interface {
	Now() time.Time
}

type SettableClock interface {
	Clock
	SetTime(t time.Time)
}

type Real struct{}

// Now implements Clock interface.
func (r Real) Now() time.Time {
	// We don't need microseconds so we truncate it to seconds.
	// It's enough for use as `created at` fields and so on.
	// TODO: should do really truncate here?
	return time.Now().Truncate(1 * time.Second).UTC()
}
