// Package debounce provides allocation-free, time-based key debouncing.
package debounce

// State tracks one physical switch. Observe may report a pending transition;
// callers must Commit it only after their event buffer has room for the event.
type State struct {
	stable      bool
	candidate   bool
	initialized bool
	since       uint32
}

// Observe advances the candidate state. It returns a transition after raw has
// remained unchanged for duration ticks. Unsigned subtraction makes timer wrap
// safe as long as duration is less than half the uint32 range.
func (s *State) Observe(raw bool, now, duration uint32) (changed, pressed bool) {
	if !s.initialized {
		s.initialized = true
		s.candidate = raw
		s.since = now
	} else if raw != s.candidate {
		s.candidate = raw
		s.since = now
	}

	if s.stable == s.candidate || now-s.since < duration {
		return false, false
	}
	return true, s.candidate
}

// Commit records a transition that has been emitted successfully.
func (s *State) Commit(pressed bool) {
	s.stable = pressed
}

// Stable returns the debounced physical state.
func (s *State) Stable() bool {
	return s.stable
}
