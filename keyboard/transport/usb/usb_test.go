package usb_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tgk-project/tgk/keyboard/report"
	usbtransport "github.com/tgk-project/tgk/keyboard/transport/usb"
)

type fakeEndpoint struct {
	ready    bool
	keyboard [][]byte
	consumer [][]byte
	raw      [][]byte
}

func (e *fakeEndpoint) Ready() bool { return e.ready }

func (e *fakeEndpoint) SendKeyboard(value []byte) bool {
	e.keyboard = append(e.keyboard, append([]byte(nil), value...))
	return true
}

func (e *fakeEndpoint) SendConsumer(value []byte) bool {
	e.consumer = append(e.consumer, append([]byte(nil), value...))
	return true
}

func (e *fakeEndpoint) SendRaw(value []byte) bool {
	e.raw = append(e.raw, append([]byte(nil), value...))
	return true
}

func newDevice(t *testing.T, endpoint *fakeEndpoint, capacity int) *usbtransport.Device {
	t.Helper()
	value, err := usbtransport.New(endpoint, capacity)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestKeyboardAndConsumerReportsUseSeparateEndpoints(t *testing.T) {
	endpoint := &fakeEndpoint{ready: true}
	value := newDevice(t, endpoint, 2)
	report := report.KeyboardReport{Modifiers: 0x02, Keys: [6]uint8{4, 5}, Consumer: 0x00e9}
	if err := value.SendKeyboardReport(report); err != nil {
		t.Fatal(err)
	}
	if got, want := endpoint.keyboard, [][]byte{{0x02, 0x00, 4, 5, 0, 0, 0, 0}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("keyboard packets = %#v, want %#v", got, want)
	}
	if got, want := endpoint.consumer, [][]byte{{0xe9, 0x00}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("consumer packets = %#v, want %#v", got, want)
	}
	if len(endpoint.raw) != 0 {
		t.Fatalf("RAW packets = %#v, want none", endpoint.raw)
	}
}

func TestRawHIDIsBoundedAndDoesNotTouchKeyboardEndpoints(t *testing.T) {
	endpoint := &fakeEndpoint{ready: true}
	value := newDevice(t, endpoint, 1)
	packet := make([]byte, usbtransport.RawReportSize)
	packet[0] = 0x7a
	if err := value.SendRaw(packet); err != nil {
		t.Fatal(err)
	}
	if len(endpoint.keyboard) != 0 || len(endpoint.consumer) != 0 {
		t.Fatal("RAW send touched keyboard or consumer endpoint")
	}
	if got := endpoint.raw[0]; !reflect.DeepEqual(got, packet) {
		t.Fatalf("RAW packet = %#v, want %#v", got, packet)
	}

	if !value.HandleRawOutput(packet) || value.HandleRawOutput(packet) {
		t.Fatal("RAW callback queue did not enforce capacity")
	}
	if got := value.RawOverflowCount(); got != 1 {
		t.Fatalf("RAW overflow count = %d, want 1", got)
	}
	dst := make([]byte, usbtransport.RawReportSize)
	if n, ok := value.ReceiveRaw(dst); !ok || n != usbtransport.RawReportSize || !reflect.DeepEqual(dst, packet) {
		t.Fatalf("ReceiveRaw() = %d, %v, %#v", n, ok, dst)
	}
}

func TestRejectsMalformedOrUnavailablePackets(t *testing.T) {
	endpoint := &fakeEndpoint{}
	value := newDevice(t, endpoint, 1)
	if err := value.SendKeyboardReport(report.KeyboardReport{}); !errors.Is(err, usbtransport.ErrEndpointNotReady) {
		t.Fatalf("SendKeyboardReport() error = %v", err)
	}
	if err := value.SendRaw(make([]byte, usbtransport.RawReportSize-1)); !errors.Is(err, usbtransport.ErrInvalidRawReport) {
		t.Fatalf("SendRaw() error = %v", err)
	}
	if value.HandleRawOutput(make([]byte, usbtransport.RawReportSize-1)) {
		t.Fatal("malformed RAW packet was accepted")
	}
	if n, ok := value.ReceiveRaw(make([]byte, usbtransport.RawReportSize-1)); ok || n != 0 {
		t.Fatalf("short ReceiveRaw() = %d, %v", n, ok)
	}
}

func TestReleaseAllClearsBothReports(t *testing.T) {
	endpoint := &fakeEndpoint{ready: true}
	value := newDevice(t, endpoint, 1)
	if err := value.SendKeyboardReport(report.KeyboardReport{Modifiers: 1, Keys: [6]uint8{4}, Consumer: 0xe9}); err != nil {
		t.Fatal(err)
	}
	if err := value.ReleaseAll(); err != nil {
		t.Fatal(err)
	}
	if got, want := endpoint.keyboard[len(endpoint.keyboard)-1], []byte{0, 0, 0, 0, 0, 0, 0, 0}; !reflect.DeepEqual(got, want) {
		t.Fatalf("released keyboard packet = %#v", got)
	}
	if got, want := endpoint.consumer[len(endpoint.consumer)-1], []byte{0, 0}; !reflect.DeepEqual(got, want) {
		t.Fatalf("released consumer packet = %#v", got)
	}
}
