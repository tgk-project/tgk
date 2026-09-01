// Package matrix scans a digital switch matrix through board-supplied pin
// adapters. It has no dependency on a concrete GPIO package.
package matrix

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/debounce"
)

var (
	ErrEmptyRows        = errors.New("matrix requires at least one row")
	ErrEmptyColumns     = errors.New("matrix requires at least one column")
	ErrInvalidPositions = errors.New("matrix position mapping does not match dimensions")
	ErrEventBufferFull  = errors.New("matrix event buffer is full")
)

// OutputPin is the GPIO boundary for a matrix row.
type OutputPin interface {
	Set(high bool)
}

// InputPin is the GPIO boundary for a matrix column.
type InputPin interface {
	Get() bool
}

// Config contains only scanning policy and pin adapters. Positions maps each
// row-major switch to a physical keyboard position; an empty map is row-major.
type Config struct {
	Rows      []OutputPin
	Columns   []InputPin
	Positions []uint16
	Source    event.Source
	ActiveLow bool
	Debounce  uint32
}

// Scanner implements event.Scanner for a digital key matrix.
type Scanner struct {
	rows      []OutputPin
	columns   []InputPin
	positions []uint16
	states    []debounce.State
	source    event.Source
	active    bool
	inactive  bool
	debounce  uint32
}

// New validates a matrix and initializes every row to its inactive level.
func New(config Config) (*Scanner, error) {
	if len(config.Rows) == 0 {
		return nil, ErrEmptyRows
	}
	if len(config.Columns) == 0 {
		return nil, ErrEmptyColumns
	}
	count := len(config.Rows) * len(config.Columns)
	if len(config.Positions) != 0 && len(config.Positions) != count {
		return nil, ErrInvalidPositions
	}
	positions := make([]uint16, count)
	if len(config.Positions) == 0 {
		for index := range positions {
			positions[index] = uint16(index)
		}
	} else {
		copy(positions, config.Positions)
	}

	active := !config.ActiveLow
	scanner := &Scanner{
		rows:      config.Rows,
		columns:   config.Columns,
		positions: positions,
		states:    make([]debounce.State, count),
		source:    config.Source,
		active:    active,
		inactive:  !active,
		debounce:  config.Debounce,
	}
	for _, row := range scanner.rows {
		row.Set(scanner.inactive)
	}
	return scanner, nil
}

// Scan samples all matrix switches and appends debounced position events to
// dst. It never appends beyond dst's capacity or retains the caller buffer.
func (s *Scanner) Scan(now uint32, dst []event.KeyEvent) ([]event.KeyEvent, error) {
	events := dst
	for rowIndex, row := range s.rows {
		row.Set(s.active)
		for columnIndex, column := range s.columns {
			index := rowIndex*len(s.columns) + columnIndex
			pressed := column.Get() == s.active
			changed, state := s.states[index].Observe(pressed, now, s.debounce)
			if !changed {
				continue
			}
			if len(events) == cap(events) {
				row.Set(s.inactive)
				return events, ErrEventBufferFull
			}
			s.states[index].Commit(state)
			events = append(events, event.KeyEvent{
				Source:    s.source,
				Position:  s.positions[index],
				Pressed:   state,
				Timestamp: now,
			})
		}
		row.Set(s.inactive)
	}
	return events, nil
}
