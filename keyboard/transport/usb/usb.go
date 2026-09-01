// Package usb adapts resolved keyboard reports to independent USB HID endpoints.
//
// It deliberately contains no TinyGo or machine imports. Board packages own
// descriptor registration and implement Endpoint, while this package owns
// report encoding and the bounded RAW HID hand-off used by the config service.
package usb

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/report"
)

const (
	// KeyboardReportSize is the boot-keyboard input report length.
	KeyboardReportSize = 8
	// ConsumerReportSize is the Consumer Control input report length.
	ConsumerReportSize = 2
	// RawReportSize is the fixed payload size of the vendor-defined RAW HID
	// interface. It excludes a report ID because RAW has its own HID interface.
	RawReportSize = 32
)

var (
	ErrInvalidRawQueueCapacity = errors.New("RAW HID queue capacity must be positive")
	ErrInvalidRawReport        = errors.New("RAW HID report must be exactly 32 bytes")
	ErrEndpointNotReady        = errors.New("USB HID endpoint is not ready")
	ErrPacketRejected          = errors.New("USB HID endpoint rejected packet")
)

// Endpoint is the board-owned USB boundary. Keyboard, consumer, and RAW HID
// use separate interrupt IN endpoints so configurator traffic cannot be
// mistaken for keyboard input. Send must consume packet before it returns.
type Endpoint interface {
	Ready() bool
	SendKeyboard(packet []byte) bool
	SendConsumer(packet []byte) bool
	SendRaw(packet []byte) bool
}

// Device implements transport.HostTransport and provides only byte-level RAW
// HID delivery. Command parsing intentionally belongs to Issue #11.
type Device struct {
	endpoint Endpoint
	raw      rawQueue

	keyboard [KeyboardReportSize]byte
	consumer [ConsumerReportSize]byte
	rawOut   [RawReportSize]byte
}

// New creates a transport with all callback-side RAW storage preallocated.
func New(endpoint Endpoint, rawQueueCapacity int) (*Device, error) {
	if endpoint == nil {
		return nil, ErrEndpointNotReady
	}
	if rawQueueCapacity <= 0 {
		return nil, ErrInvalidRawQueueCapacity
	}
	return &Device{
		endpoint: endpoint,
		raw: rawQueue{
			reports: make([]rawReport, rawQueueCapacity),
		},
	}, nil
}

// Ready reports whether the board USB adapter has completed endpoint setup.
func (d *Device) Ready() bool { return d.endpoint.Ready() }

// SendKeyboardReport sends the boot keyboard and Consumer Control views of
// the Primary-owned report state over their independent endpoints.
func (d *Device) SendKeyboardReport(value report.KeyboardReport) error {
	if !d.Ready() {
		return ErrEndpointNotReady
	}
	d.keyboard[0] = value.Modifiers
	d.keyboard[1] = 0
	copy(d.keyboard[2:], value.Keys[:])
	if !d.endpoint.SendKeyboard(d.keyboard[:]) {
		return ErrPacketRejected
	}

	d.consumer[0] = byte(value.Consumer)
	d.consumer[1] = byte(value.Consumer >> 8)
	if !d.endpoint.SendConsumer(d.consumer[:]) {
		return ErrPacketRejected
	}
	return nil
}

// ReleaseAll clears both Host-visible input reports. It is safe to call after
// a queue overflow, split disconnect, or Host transport reset.
func (d *Device) ReleaseAll() error {
	if !d.Ready() {
		return ErrEndpointNotReady
	}
	clear(d.keyboard[:])
	clear(d.consumer[:])
	if !d.endpoint.SendKeyboard(d.keyboard[:]) {
		return ErrPacketRejected
	}
	if !d.endpoint.SendConsumer(d.consumer[:]) {
		return ErrPacketRejected
	}
	return nil
}

// SendRaw writes one fixed-size vendor-defined report. It does not interpret
// its contents and cannot affect the keyboard or consumer endpoints.
func (d *Device) SendRaw(value []byte) error {
	if len(value) != RawReportSize {
		return ErrInvalidRawReport
	}
	if !d.Ready() {
		return ErrEndpointNotReady
	}
	copy(d.rawOut[:], value)
	if !d.endpoint.SendRaw(d.rawOut[:]) {
		return ErrPacketRejected
	}
	return nil
}

// HandleRawOutput is for the board's USB callback. It only copies a valid
// fixed-size packet into a bounded queue; the Primary event loop must later
// call ReceiveRaw to run any configuration work outside the callback.
func (d *Device) HandleRawOutput(value []byte) bool {
	return d.raw.push(value)
}

// ReceiveRaw copies the next RAW HID packet into dst. A destination shorter
// than RawReportSize is rejected without removing the queued packet.
func (d *Device) ReceiveRaw(dst []byte) (int, bool) {
	return d.raw.pop(dst)
}

// RawOverflowCount reports dropped callback packets without coupling this
// transport package to a logging implementation.
func (d *Device) RawOverflowCount() uint32 { return d.raw.overflowCount() }
