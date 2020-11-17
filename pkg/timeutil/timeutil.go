package timeutil

import "time"

// Now now struct
type Now struct {
	time.Time
}

// New initialize Now with time
func New(t time.Time) *Now {
	return &Now{t}
}

// BeginningOfYear beginning of year
func BeginningOfYear() time.Time {
	return New(time.Now()).BeginningOfYear()
}

// BeginningOfYear BeginningOfYear beginning of year
func (now *Now) BeginningOfYear() time.Time {
	y, _, _ := now.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, now.Location())
}
