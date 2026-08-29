// Package topre scans analogue Topre-style sensors through a board-provided
// reader. Keycodes and keyboard configuration remain outside this package.
package topre

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/debounce"
)

var (
	ErrNoSensors         = errors.New("topre requires at least one sensor")
	ErrInvalidThresholds = errors.New("topre release threshold must be below press threshold")
	ErrEventBufferFull   = errors.New("topre event buffer is full")
)

// Reader returns one sensor value. The index is scanner-local; Config.Positions
// maps it to the physical Position emitted to the keyboard core.
type Reader interface {
	Read(index uint16) (uint16, error)
}

// Config describes Topre sensor mapping and hysteresis.
type Config struct {
	Positions        []uint16
	Source           event.Source
	PressThreshold   uint16
	ReleaseThreshold uint16
	Debounce         uint32
}

// Scanner implements event.Scanner for Topre-style analogue sensors.
type Scanner struct {
	reader    Reader
	positions []uint16
	states    []debounce.State
	source    event.Source
	press     uint16
	release   uint16
	debounce  uint32
}

// New validates a Topre scanner and preallocates all per-sensor state.
func New(config Config, reader Reader) (*Scanner, error) {
	if len(config.Positions) == 0 {
		return nil, ErrNoSensors
	}
	if config.ReleaseThreshold >= config.PressThreshold {
		return nil, ErrInvalidThresholds
	}
	positions := make([]uint16, len(config.Positions))
	copy(positions, config.Positions)
	return &Scanner{
		reader:    reader,
		positions: positions,
		states:    make([]debounce.State, len(positions)),
		source:    config.Source,
		press:     config.PressThreshold,
		release:   config.ReleaseThreshold,
		debounce:  config.Debounce,
	}, nil
}

// Scan reads each sensor and appends only debounced physical Position events.
func (s *Scanner) Scan(now uint32, dst []event.KeyEvent) ([]event.KeyEvent, error) {
	events := dst
	for index := range s.positions {
		value, err := s.reader.Read(uint16(index))
		if err != nil {
			return events, err
		}
		pressed := value >= s.press
		if s.states[index].Stable() {
			pressed = value > s.release
		}
		changed, state := s.states[index].Observe(pressed, now, s.debounce)
		if !changed {
			continue
		}
		if len(events) == cap(events) {
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
	return events, nil
}
