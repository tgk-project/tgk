// Package transport declares the boundaries between the core and I/O adapters.
package transport

import (
	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/keyboard/report"
)

// HostTransport sends resolved reports to the Host. Implementations may be USB
// or BLE, but neither technology is visible to the keymap engine.
type HostTransport interface {
	Ready() bool
	SendKeyboardReport(report.KeyboardReport) error
	ReleaseAll() error
}

// SplitTransport exchanges only physical position events. It never carries
// keycodes, layers, modifier state, or Host reports.
type SplitTransport interface {
	SendEvents([]event.KeyEvent) error
	ReceiveEvents(dst []event.KeyEvent) ([]event.KeyEvent, error)
}
