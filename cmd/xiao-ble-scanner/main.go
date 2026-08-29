//go:build xiao_ble

// xiao-ble-scanner is a build probe for the XIAO BLE matrix adapter. Its 1x1
// wiring is intentionally only a minimal board bring-up configuration; a
// keyboard project supplies its own matrix wiring and position map.
package main

import (
	"machine"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/platform/xiao_ble"
)

func main() {
	scanner, err := xiaoble.NewMatrix(xiaoble.MatrixConfig{
		Rows:      []machine.Pin{machine.D0},
		Columns:   []machine.Pin{machine.D1},
		Positions: []uint16{0},
		Source:    0,
		ActiveLow: true,
	})
	if err != nil {
		return
	}
	var events [1]event.KeyEvent
	_, _ = scanner.Scan(0, events[:0])
}
