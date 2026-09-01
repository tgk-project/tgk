//go:build xiao_ble

// Package xiaoble adapts the TinyGo XIAO BLE machine package to the generic
// matrix scanner. Board wiring remains an explicit caller-provided value.
package xiaoble

import (
	"machine"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/matrix"
)

// MatrixConfig is the XIAO BLE-specific wiring input for a matrix scanner.
// It deliberately contains physical positions, not keycodes or keymap data.
type MatrixConfig struct {
	Rows      []machine.Pin
	Columns   []machine.Pin
	Positions []uint16
	Source    event.Source
	ActiveLow bool
	Debounce  uint32
}

// NewMatrix configures XIAO BLE pins and adapts them to matrix.Scanner.
func NewMatrix(config MatrixConfig) (*matrix.Scanner, error) {
	rows := make([]matrix.OutputPin, len(config.Rows))
	for index, pin := range config.Rows {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		rows[index] = outputPin{pin: pin}
	}
	columns := make([]matrix.InputPin, len(config.Columns))
	for index, pin := range config.Columns {
		pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		columns[index] = inputPin{pin: pin}
	}
	return matrix.New(matrix.Config{
		Rows:      rows,
		Columns:   columns,
		Positions: config.Positions,
		Source:    config.Source,
		ActiveLow: config.ActiveLow,
		Debounce:  config.Debounce,
	})
}

type outputPin struct{ pin machine.Pin }

func (p outputPin) Set(high bool) { p.pin.Set(high) }

type inputPin struct{ pin machine.Pin }

func (p inputPin) Get() bool { return p.pin.Get() }
