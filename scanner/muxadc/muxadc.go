// Package muxadc scans thresholded analogue switches through multiplexer and
// ADC adapters. Board wiring stays outside this package.
package muxadc

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/debounce"
)

var (
	ErrNoRows            = errors.New("muxadc requires at least one row")
	ErrNoChannels        = errors.New("muxadc requires at least one channel")
	ErrInvalidThresholds = errors.New("muxadc release threshold must be below press threshold")
	ErrInvalidPositions  = errors.New("muxadc position mapping does not match dimensions")
	ErrEventBufferFull   = errors.New("muxadc event buffer is full")
)

// Multiplexer selects one analogue channel.
type Multiplexer interface {
	Select(channel uint8)
}

// RowDriver energizes one sensor row for an ADC sample.
type RowDriver interface {
	Drive(row uint8, high bool)
}

// ADC is the analogue input boundary. The active multiplexer channel is chosen
// by Multiplexer before Read is called.
type ADC interface {
	Read() (uint16, error)
}

// Config describes a row-by-channel sensor array and its hysteresis policy.
type Config struct {
	Rows             uint8
	Channels         []uint8
	Positions        []uint16
	Source           event.Source
	PressThreshold   uint16
	ReleaseThreshold uint16
	Debounce         uint32
}

// Scanner implements event.Scanner for MUX/ADC switch hardware.
type Scanner struct {
	mux       Multiplexer
	rows      RowDriver
	adc       ADC
	channels  []uint8
	positions []uint16
	states    []debounce.State
	source    event.Source
	press     uint16
	release   uint16
	debounce  uint32
	rowCount  uint8
}

// New validates a scanner and preallocates its state. The adapter interfaces
// can be implemented by GPIO/ADC code without leaking into the keyboard core.
func New(config Config, mux Multiplexer, rows RowDriver, adc ADC) (*Scanner, error) {
	if config.Rows == 0 {
		return nil, ErrNoRows
	}
	if len(config.Channels) == 0 {
		return nil, ErrNoChannels
	}
	if config.ReleaseThreshold >= config.PressThreshold {
		return nil, ErrInvalidThresholds
	}
	count := int(config.Rows) * len(config.Channels)
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
	return &Scanner{
		mux:       mux,
		rows:      rows,
		adc:       adc,
		channels:  config.Channels,
		positions: positions,
		states:    make([]debounce.State, count),
		source:    config.Source,
		press:     config.PressThreshold,
		release:   config.ReleaseThreshold,
		debounce:  config.Debounce,
		rowCount:  config.Rows,
	}, nil
}

// Scan samples every channel and row. Values in the hysteresis region preserve
// the previous debounced state, preventing threshold chatter.
func (s *Scanner) Scan(now uint32, dst []event.KeyEvent) ([]event.KeyEvent, error) {
	events := dst
	for channelIndex, channel := range s.channels {
		s.mux.Select(channel)
		for row := uint8(0); row < s.rowCount; row++ {
			s.rows.Drive(row, true)
			value, err := s.adc.Read()
			s.rows.Drive(row, false)
			if err != nil {
				return events, err
			}

			index := channelIndex*int(s.rowCount) + int(row)
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
	}
	return events, nil
}
