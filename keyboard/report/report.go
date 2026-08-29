// Package report contains transport-neutral Host report values.
package report

// ErrorRollOver is the HID keyboard usage reported in every key slot when
// more than six ordinary keys are active in a boot-style report.
const ErrorRollOver uint8 = 0x01

// KeyboardReport is the state required by a standard keyboard input report and
// a consumer-control report. Transport implementations decide how to encode it.
type KeyboardReport struct {
	Modifiers uint8
	Keys      [6]uint8
	Consumer  uint16
}
