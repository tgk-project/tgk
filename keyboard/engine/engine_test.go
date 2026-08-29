package engine_test

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/tgk-project/tgk/keyboard/engine"
	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/keyboard/keymap"
	"github.com/tgk-project/tgk/keyboard/report"
)

const (
	keyA = 4
	keyB = 5
	keyC = 6
)

type fakeHost struct {
	ready    bool
	reports  []report.KeyboardReport
	releases int
}

func (h *fakeHost) Ready() bool { return h.ready }

func (h *fakeHost) SendKeyboardReport(value report.KeyboardReport) error {
	h.reports = append(h.reports, value)
	return nil
}

func (h *fakeHost) ReleaseAll() error {
	h.releases++
	return nil
}

type fakeScanner struct{ events []event.KeyEvent }

func (s fakeScanner) Scan(_ uint32, dst []event.KeyEvent) ([]event.KeyEvent, error) {
	return append(dst, s.events...), nil
}

type fakeSplit struct{ events []event.KeyEvent }

func (s fakeSplit) SendEvents([]event.KeyEvent) error { return nil }

func (s fakeSplit) ReceiveEvents(dst []event.KeyEvent) ([]event.KeyEvent, error) {
	return append(dst, s.events...), nil
}

func testKeymap(t *testing.T) keymap.Keymap {
	t.Helper()
	km, err := keymap.New([][]keymap.Binding{
		{
			{Behavior: keymap.BehaviorModifier, Param1: 0x02}, // position 0
			{Behavior: keymap.BehaviorKey, Param1: keyA},
			{Behavior: keymap.BehaviorConsumer, Param1: 0x00e9},
			{Behavior: keymap.BehaviorMomentaryLayer, Param1: 1},
			{Behavior: keymap.BehaviorToLayer, Param1: 2},
			{Behavior: keymap.BehaviorToggleLayer, Param1: 1},
		},
		{
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorKey, Param1: keyB},
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorNone},
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorTransparent},
		},
		{
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorKey, Param1: keyC},
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorTransparent},
			{Behavior: keymap.BehaviorToLayer, Param1: 0},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return km
}

func newEngine(t *testing.T, host *fakeHost) *engine.Engine {
	t.Helper()
	value, err := engine.New(testKeymap(t), host, 16)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func runEvents(t *testing.T, value *engine.Engine, events ...event.KeyEvent) {
	t.Helper()
	for _, item := range events {
		if !value.Enqueue(item) {
			t.Fatal("Enqueue() unexpectedly returned false")
		}
	}
	if _, err := value.Run(0); err != nil {
		t.Fatal(err)
	}
}

func keyEvent(position uint16, pressed bool) event.KeyEvent {
	return event.KeyEvent{Source: 1, Position: position, Pressed: pressed}
}

func TestPressReleaseOrderingAndReport(t *testing.T) {
	host := &fakeHost{ready: true}
	value := newEngine(t, host)

	runEvents(t, value,
		keyEvent(0, true),
		keyEvent(1, true),
		keyEvent(1, false),
		keyEvent(0, false),
	)

	want := []report.KeyboardReport{
		{Modifiers: 0x02},
		{Modifiers: 0x02, Keys: [6]uint8{keyA}},
		{Modifiers: 0x02},
		{},
	}
	if !reflect.DeepEqual(host.reports, want) {
		t.Fatalf("reports = %#v, want %#v", host.reports, want)
	}
}

func TestLayerResolutionTransparentMOAndTG(t *testing.T) {
	host := &fakeHost{ready: true}
	value := newEngine(t, host)

	runEvents(t, value,
		keyEvent(3, true), // MO(1)
		keyEvent(0, true), // transparent layer 1 falls through to modifier
		keyEvent(0, false),
		keyEvent(1, true),  // layer 1 B
		keyEvent(3, false), // release MO while B remains held
		keyEvent(1, false), // must release the resolved B binding, not A
		keyEvent(5, true),  // TG(1) on
		keyEvent(5, false),
		keyEvent(1, true),
		keyEvent(1, false),
		keyEvent(5, true), // TG(1) off through transparent binding
		keyEvent(5, false),
		keyEvent(1, true),
		keyEvent(1, false),
	)

	want := []report.KeyboardReport{
		{Modifiers: 0x02}, {},
		{Keys: [6]uint8{keyB}}, {},
		{Keys: [6]uint8{keyB}}, {},
		{Keys: [6]uint8{keyA}}, {},
	}
	if !reflect.DeepEqual(host.reports, want) {
		t.Fatalf("reports = %#v, want %#v", host.reports, want)
	}
}

func TestLayerMoveAndConsumer(t *testing.T) {
	host := &fakeHost{ready: true}
	value := newEngine(t, host)

	runEvents(t, value,
		keyEvent(4, true), keyEvent(4, false), // TO(2)
		keyEvent(1, true), keyEvent(1, false), // C
		keyEvent(5, true), keyEvent(5, false), // TO(0)
		keyEvent(2, true), keyEvent(2, false), // consumer volume up
	)

	want := []report.KeyboardReport{
		{Keys: [6]uint8{keyC}}, {},
		{Consumer: 0x00e9}, {},
	}
	if !reflect.DeepEqual(host.reports, want) {
		t.Fatalf("reports = %#v, want %#v", host.reports, want)
	}
}

func TestScannerAndSplitInputResolveIdentically(t *testing.T) {
	events := []event.KeyEvent{keyEvent(1, true), keyEvent(1, false)}

	fromScannerHost := &fakeHost{ready: true}
	fromScanner := newEngine(t, fromScannerHost)
	count, err := fromScanner.PollScanner(fakeScanner{events: events}, 100, make([]event.KeyEvent, 0, 4))
	if err != nil || count != len(events) {
		t.Fatalf("PollScanner() = %d, %v", count, err)
	}
	if _, err := fromScanner.Run(0); err != nil {
		t.Fatal(err)
	}

	fromSplitHost := &fakeHost{ready: true}
	fromSplit := newEngine(t, fromSplitHost)
	count, err = fromSplit.ReceiveSplit(fakeSplit{events: events}, make([]event.KeyEvent, 0, 4))
	if err != nil || count != len(events) {
		t.Fatalf("ReceiveSplit() = %d, %v", count, err)
	}
	if _, err := fromSplit.Run(0); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fromScannerHost.reports, fromSplitHost.reports) {
		t.Fatalf("scanner reports = %#v, split reports = %#v", fromScannerHost.reports, fromSplitHost.reports)
	}
}

func TestSourceKeepsSamePositionIndependent(t *testing.T) {
	host := &fakeHost{ready: true}
	value := newEngine(t, host)

	runEvents(t, value,
		event.KeyEvent{Source: 1, Position: 1, Pressed: true},
		event.KeyEvent{Source: 2, Position: 1, Pressed: true},
		event.KeyEvent{Source: 1, Position: 1, Pressed: false},
		event.KeyEvent{Source: 2, Position: 1, Pressed: false},
	)

	want := []report.KeyboardReport{
		{Keys: [6]uint8{keyA}},
		{Keys: [6]uint8{keyA}},
		{Keys: [6]uint8{keyA}},
		{},
	}
	if !reflect.DeepEqual(host.reports, want) {
		t.Fatalf("reports = %#v, want %#v", host.reports, want)
	}
}

func TestSixKeyReportRecoversAfterRollover(t *testing.T) {
	layer := make([]keymap.Binding, 7)
	for i := range layer {
		layer[i] = keymap.Binding{Behavior: keymap.BehaviorKey, Param1: uint32(keyA + i)}
	}
	km, err := keymap.New([][]keymap.Binding{layer})
	if err != nil {
		t.Fatal(err)
	}
	host := &fakeHost{ready: true}
	value, err := engine.New(km, host, 16)
	if err != nil {
		t.Fatal(err)
	}

	for position := uint16(0); position < 7; position++ {
		runEvents(t, value, keyEvent(position, true))
	}
	if got, want := host.reports[len(host.reports)-1].Keys, [6]uint8{
		report.ErrorRollOver,
		report.ErrorRollOver,
		report.ErrorRollOver,
		report.ErrorRollOver,
		report.ErrorRollOver,
		report.ErrorRollOver,
	}; got != want {
		t.Fatalf("rollover keys = %#v, want %#v", got, want)
	}

	runEvents(t, value, keyEvent(6, false))
	if got, want := host.reports[len(host.reports)-1].Keys, [6]uint8{4, 5, 6, 7, 8, 9}; got != want {
		t.Fatalf("recovered keys = %#v, want %#v", got, want)
	}
}

func TestOverflowReleasesStateAndDropsUnmatchedEvents(t *testing.T) {
	host := &fakeHost{ready: true}
	value, err := engine.New(testKeymap(t), host, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !value.Enqueue(keyEvent(1, true)) {
		t.Fatal("first Enqueue() returned false")
	}
	if value.Enqueue(keyEvent(1, false)) {
		t.Fatal("second Enqueue() returned true despite full queue")
	}
	processed, err := value.Run(0)
	if processed != 0 || !errors.Is(err, engine.ErrQueueOverflow) {
		t.Fatalf("Run() = %d, %v; want 0, ErrQueueOverflow", processed, err)
	}
	if host.releases != 1 {
		t.Fatalf("ReleaseAll calls = %d, want 1", host.releases)
	}
	if got := value.Report(); got != (report.KeyboardReport{}) {
		t.Fatalf("Report() = %#v, want zero report", got)
	}
	if got := value.QueueStats().OverflowCount; got != 1 {
		t.Fatalf("OverflowCount = %d, want 1", got)
	}
}

func TestConcurrentCallbacksOnlyEnqueue(t *testing.T) {
	km, err := keymap.New([][]keymap.Binding{{{Behavior: keymap.BehaviorNone}}})
	if err != nil {
		t.Fatal(err)
	}
	value, err := engine.NewWithOptions(km, nil, engine.Options{QueueCapacity: 64, MaxPressed: 64})
	if err != nil {
		t.Fatal(err)
	}

	var group sync.WaitGroup
	for producer := 0; producer < 4; producer++ {
		group.Add(1)
		go func(offset int) {
			defer group.Done()
			for i := 0; i < 16; i++ {
				if !value.Enqueue(event.KeyEvent{Source: event.Source(offset + i), Pressed: true}) {
					t.Errorf("Enqueue() unexpectedly rejected producer event")
				}
			}
		}(producer * 16)
	}
	group.Wait()

	processed, err := value.Run(0)
	if err != nil || processed != 64 {
		t.Fatalf("Run() = %d, %v; want 64, nil", processed, err)
	}
	if got := value.QueueStats().OverflowCount; got != 0 {
		t.Fatalf("OverflowCount = %d, want 0", got)
	}
}

func TestEventPathDoesNotAllocate(t *testing.T) {
	km, err := keymap.New([][]keymap.Binding{{{Behavior: keymap.BehaviorKey, Param1: keyA}}})
	if err != nil {
		t.Fatal(err)
	}
	value, err := engine.New(km, nil, 2)
	if err != nil {
		t.Fatal(err)
	}
	down := keyEvent(0, true)
	up := keyEvent(0, false)

	allocations := testing.AllocsPerRun(1_000, func() {
		if !value.Enqueue(down) {
			t.Fatal("enqueue down")
		}
		if _, err := value.Run(1); err != nil {
			t.Fatal(err)
		}
		if !value.Enqueue(up) {
			t.Fatal("enqueue up")
		}
		if _, err := value.Run(1); err != nil {
			t.Fatal(err)
		}
	})
	if allocations != 0 {
		t.Fatalf("event path allocations = %f, want 0", allocations)
	}
}
