package muxadc_test

import (
	"testing"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/muxadc"
)

type multiplexer struct{ selected uint8 }

func (m *multiplexer) Select(channel uint8) { m.selected = channel }

type rowDriver struct {
	active uint8
	high   bool
}

func (d *rowDriver) Drive(row uint8, high bool) {
	d.active = row
	d.high = high
}

type adc struct {
	mux    *multiplexer
	driver *rowDriver
	values map[[2]uint8]uint16
}

func (a *adc) Read() (uint16, error) {
	return a.values[[2]uint8{a.mux.selected, a.driver.active}], nil
}

func TestMUXADCDebounceHysteresisAndPositionMapping(t *testing.T) {
	mux := &multiplexer{}
	driver := &rowDriver{}
	reader := &adc{mux: mux, driver: driver, values: map[[2]uint8]uint16{{3, 0}: 900}}
	scanner, err := muxadc.New(muxadc.Config{
		Rows:             1,
		Channels:         []uint8{3},
		Positions:        []uint16{71},
		Source:           5,
		PressThreshold:   800,
		ReleaseThreshold: 600,
		Debounce:         3,
	}, mux, driver, reader)
	if err != nil {
		t.Fatal(err)
	}
	buffer := make([]event.KeyEvent, 0, 1)
	if events, err := scanner.Scan(0, buffer); err != nil || len(events) != 0 {
		t.Fatalf("first Scan() = %#v, %v; want no event", events, err)
	}
	events, err := scanner.Scan(3, buffer)
	if err != nil {
		t.Fatal(err)
	}
	wantDown := event.KeyEvent{Source: 5, Position: 71, Pressed: true, Timestamp: 3}
	if len(events) != 1 || events[0] != wantDown {
		t.Fatalf("press event = %#v, want %#v", events, wantDown)
	}

	reader.values[[2]uint8{3, 0}] = 700 // Within hysteresis: remains pressed.
	if events, err := scanner.Scan(4, buffer); err != nil || len(events) != 0 {
		t.Fatalf("hysteresis Scan() = %#v, %v; want no event", events, err)
	}
	reader.values[[2]uint8{3, 0}] = 500
	if events, err := scanner.Scan(5, buffer); err != nil || len(events) != 0 {
		t.Fatalf("release debounce start = %#v, %v; want no event", events, err)
	}
	events, err = scanner.Scan(8, buffer)
	if err != nil {
		t.Fatal(err)
	}
	wantUp := event.KeyEvent{Source: 5, Position: 71, Timestamp: 8}
	if len(events) != 1 || events[0] != wantUp {
		t.Fatalf("release event = %#v, want %#v", events, wantUp)
	}
	if mux.selected != 3 || driver.high {
		t.Fatalf("mux/row state = channel %d, high %t; expected selected channel and inactive row", mux.selected, driver.high)
	}
}
