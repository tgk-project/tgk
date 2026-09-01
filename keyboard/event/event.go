// Package event defines physical keyboard input shared by local and split sources.
package event

// Source identifies the physical producer of a key event. It is intentionally
// independent of a transport connection or board implementation.
type Source uint8

// KeyEvent describes one physical position transition. A Position is never a
// keycode: keymap resolution happens only in the Primary engine.
type KeyEvent struct {
	Source    Source
	Position  uint16
	Pressed   bool
	Sequence  uint16
	Timestamp uint32
}

// Scanner emits physical position events into a caller-owned buffer. A scanner
// must append to dst[:0] and must not retain dst after Scan returns.
type Scanner interface {
	Scan(now uint32, dst []KeyEvent) ([]KeyEvent, error)
}
