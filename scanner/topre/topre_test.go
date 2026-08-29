package topre_test

import (
	"testing"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/topre"
)

type reader struct{ values []uint16 }

func (r *reader) Read(index uint16) (uint16, error) {
	return r.values[index], nil
}

func TestTopreDebounceHysteresisAndPositionMapping(t *testing.T) {
	input := &reader{values: []uint16{900}}
	scanner, err := topre.New(topre.Config{
		Positions:        []uint16{91},
		Source:           7,
		PressThreshold:   800,
		ReleaseThreshold: 600,
		Debounce:         2,
	}, input)
	if err != nil {
		t.Fatal(err)
	}
	buffer := make([]event.KeyEvent, 0, 1)
	if events, err := scanner.Scan(10, buffer); err != nil || len(events) != 0 {
		t.Fatalf("first Scan() = %#v, %v; want no event", events, err)
	}
	events, err := scanner.Scan(12, buffer)
	if err != nil {
		t.Fatal(err)
	}
	wantDown := event.KeyEvent{Source: 7, Position: 91, Pressed: true, Timestamp: 12}
	if len(events) != 1 || events[0] != wantDown {
		t.Fatalf("press event = %#v, want %#v", events, wantDown)
	}

	input.values[0] = 700
	if events, err := scanner.Scan(13, buffer); err != nil || len(events) != 0 {
		t.Fatalf("hysteresis Scan() = %#v, %v; want no event", events, err)
	}
	input.values[0] = 500
	if events, err := scanner.Scan(14, buffer); err != nil || len(events) != 0 {
		t.Fatalf("release debounce start = %#v, %v; want no event", events, err)
	}
	events, err = scanner.Scan(16, buffer)
	if err != nil {
		t.Fatal(err)
	}
	wantUp := event.KeyEvent{Source: 7, Position: 91, Timestamp: 16}
	if len(events) != 1 || events[0] != wantUp {
		t.Fatalf("release event = %#v, want %#v", events, wantUp)
	}
}
