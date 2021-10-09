package monday

import "sync"

var keyExists = struct{}{}

// Set is a thread safe set data structure.
//
// It is ported from https://github.com/fatih/set with only the required functionality.
type set struct {
	m map[Locale]struct{}
	l sync.RWMutex
}

// NewSet allocates and returns a new Set. It accepts a variable number of
// arguments to populate the initial set. If nothing is passed a Set with zero
// size is created.
func newSet(items ...Locale) *set {
	s := set{
		m: make(map[Locale]struct{}),
	}
	s.Add(items...)
	return &s
}

// Add adds the specified items (one or more) to the set. The underlying
// Set s is modified. If no items are passed it silently returns.
func (s *set) Add(items ...Locale) {
	if len(items) == 0 {
		return
	}
	s.l.Lock()
	defer s.l.Unlock()
	for _, item := range items {
		s.m[item] = keyExists
	}
}

// Each traverses the items in the Set, calling the provided function f for
// each set member. Traversal will continue until all items in the Set have
// been visited, or if the closure returns false.
func (s *set) Each(f func(item Locale) bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	for item := range s.m {
		if !f(item) {
			break
		}
	}
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of the items exist.
func (s *set) Has(items ...Locale) bool {
	if len(items) == 0 {
		return false
	}
	s.l.RLock()
	defer s.l.RUnlock()
	has := true
	for _, item := range items {
		if _, has = s.m[item]; !has {
			break
		}
	}
	return has
}
