// Package via adapts the VIA RAW HID keymap subset to config.Service. It owns
// command IDs and byte encoding; the service itself has no VIA or USB imports.
package via

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/config/service"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

const (
	ReportSize      = 32
	ProtocolVersion = 9

	GetProtocolVersion byte = 0x01
	GetKeyboardValue   byte = 0x02
	SetKeyboardValue   byte = 0x03
	GetKeycode         byte = 0x04
	SetKeycode         byte = 0x05
	EEPROMReset        byte = 0x0a
	GetLayerCount      byte = 0x11
	GetKeymapBuffer    byte = 0x12
	SetKeymapBuffer    byte = 0x13

	keyboardValueFirmwareVersion byte = 0x04
)

var (
	ErrInvalidPacket = errors.New("VIA RAW HID packet must be 32 bytes")
	ErrInvalidMatrix = errors.New("VIA matrix dimensions do not match config service")
	ErrUnsupported   = errors.New("VIA command or keycode is unsupported")
)

// Matrix maps VIA's row/column addressing to the dense physical positions
// owned by Config Service. A keyboard definition supplies the dimensions.
type Matrix struct {
	Rows    uint8
	Columns uint8
}

// Adapter is invoked by the Primary event loop after USB has copied a packet
// out of its bounded callback queue.
type Adapter struct {
	service *service.Service
	matrix  Matrix
}

// RawTransport is satisfied by keyboard/transport/usb.Device. Keeping the
// byte boundary here avoids coupling the parser to USB endpoint details.
type RawTransport interface {
	ReceiveRaw(dst []byte) (int, bool)
	SendRaw(value []byte) error
}

// Runner owns one preallocated request buffer and processes at most one RAW
// HID command per event-loop turn. USB callbacks only enqueue packets; they
// do not call Handle, keymap replacement, or persistence.
type Runner struct {
	adapter *Adapter
	raw     RawTransport
	request [ReportSize]byte
}

func New(value *service.Service, matrix Matrix) (*Adapter, error) {
	if value == nil || matrix.Rows == 0 || matrix.Columns == 0 || int(matrix.Rows)*int(matrix.Columns) != int(value.PositionCount()) {
		return nil, ErrInvalidMatrix
	}
	return &Adapter{service: value, matrix: matrix}, nil
}

func NewRunner(adapter *Adapter, raw RawTransport) (*Runner, error) {
	if adapter == nil || raw == nil {
		return nil, ErrInvalidPacket
	}
	return &Runner{adapter: adapter, raw: raw}, nil
}

// ProcessNext handles one queued packet and sends a fixed-size response. A
// parser error is returned after its zero-filled response has been sent so the
// event loop can count or log rejected input without leaving a browser request
// indefinitely pending.
func (r *Runner) ProcessNext(now uint32) (bool, error) {
	n, ok := r.raw.ReceiveRaw(r.request[:])
	if !ok {
		return false, nil
	}
	if n != ReportSize {
		return true, ErrInvalidPacket
	}
	response, handleErr := r.adapter.Handle(r.request[:], now)
	if err := r.raw.SendRaw(response[:]); err != nil {
		return true, err
	}
	return true, handleErr
}

// Handle returns a same-size response with the request command echoed at byte
// zero. On a rejected command the response remains zero-filled and the error
// is returned for composition-level observability; no malformed input panics.
func (a *Adapter) Handle(packet []byte, now uint32) ([ReportSize]byte, error) {
	var response [ReportSize]byte
	if len(packet) != ReportSize {
		return response, ErrInvalidPacket
	}
	response[0] = packet[0]
	switch packet[0] {
	case GetProtocolVersion:
		response[1] = byte(ProtocolVersion >> 8)
		response[2] = byte(ProtocolVersion)
		return response, nil
	case GetKeyboardValue:
		if packet[1] != keyboardValueFirmwareVersion {
			return response, ErrUnsupported
		}
		// The immutable VID/PID and product name are selected by VIA through
		// the keyboard definition JSON. A stable non-zero firmware version is
		// returned on the standard VIA keyboard-value path.
		response[1] = packet[1]
		response[2], response[3], response[4], response[5] = 'T', 'G', 'K', '2'
		return response, nil
	case GetLayerCount:
		response[1] = a.service.LayerCount()
		return response, nil
	case GetKeycode:
		binding, err := a.binding(packet[1], packet[2], packet[3])
		if err != nil {
			return response, err
		}
		code, err := encode(binding)
		if err != nil {
			return response, err
		}
		response[1], response[2], response[3] = packet[1], packet[2], packet[3]
		response[4], response[5] = byte(code>>8), byte(code)
		return response, nil
	case SetKeycode:
		binding, err := decode(uint16(packet[4])<<8 | uint16(packet[5]))
		if err != nil {
			return response, err
		}
		position, err := a.positionFor(packet[2], packet[3])
		if err != nil {
			return response, err
		}
		if err := a.service.SetBinding(now, packet[1], position, binding); err != nil {
			return response, err
		}
		response[1], response[2], response[3], response[4], response[5] = packet[1], packet[2], packet[3], packet[4], packet[5]
		return response, nil
	case GetKeymapBuffer:
		offset, size := int(packet[1])<<8|int(packet[2]), int(packet[3])
		data, err := a.buffer(offset, size)
		if err != nil {
			return response, err
		}
		response[1], response[2], response[3] = packet[1], packet[2], packet[3]
		copy(response[4:], data)
		return response, nil
	case SetKeymapBuffer:
		offset, size := int(packet[1])<<8|int(packet[2]), int(packet[3])
		if size > ReportSize-4 || size%2 != 0 || offset%2 != 0 || offset+size > a.bufferSize() {
			return response, service.ErrInvalidBuffer
		}
		next, err := a.allBindings()
		if err != nil {
			return response, err
		}
		for index := 0; index < size; index += 2 {
			binding, err := decode(uint16(packet[4+index])<<8 | uint16(packet[5+index]))
			if err != nil {
				return response, err
			}
			next[(offset+index)/2] = binding
		}
		if err := a.service.ReplaceBindings(now, next); err != nil {
			return response, err
		}
		response[1], response[2], response[3] = packet[1], packet[2], packet[3]
		return response, nil
	case EEPROMReset:
		if err := a.service.FactoryReset(now); err != nil {
			return response, err
		}
		return response, nil
	default:
		return response, ErrUnsupported
	}
}

// Commit persists a staged config at an event-loop controlled safe point.
func (a *Adapter) Commit() error { return a.service.Commit() }

func (a *Adapter) binding(layer, row, column byte) (keymap.Binding, error) {
	position, err := a.positionFor(row, column)
	if err != nil {
		return keymap.Binding{}, err
	}
	value, ok := a.service.Binding(layer, position)
	if !ok {
		return keymap.Binding{}, service.ErrInvalidAddress
	}
	return value, nil
}

func (a *Adapter) position(row, column byte) uint16 {
	return uint16(row)*uint16(a.matrix.Columns) + uint16(column)
}

func (a *Adapter) positionFor(row, column byte) (uint16, error) {
	if row >= a.matrix.Rows || column >= a.matrix.Columns {
		return 0, service.ErrInvalidAddress
	}
	return a.position(row, column), nil
}

func (a *Adapter) buffer(offset, size int) ([]byte, error) {
	if size > ReportSize-4 || offset%2 != 0 || size%2 != 0 || offset < 0 || size < 0 || offset+size > a.bufferSize() {
		return nil, service.ErrInvalidBuffer
	}
	bindings, err := a.allBindings()
	if err != nil {
		return nil, err
	}
	encoded := make([]byte, a.bufferSize())
	for index, binding := range bindings {
		code, err := encode(binding)
		if err != nil {
			return nil, err
		}
		encoded[index*2], encoded[index*2+1] = byte(code>>8), byte(code)
	}
	return append([]byte(nil), encoded[offset:offset+size]...), nil
}

func (a *Adapter) allBindings() ([]keymap.Binding, error) {
	return a.service.Bindings(), nil
}

func (a *Adapter) bufferSize() int {
	return int(a.service.LayerCount()) * int(a.service.PositionCount()) * 2
}

// VIA's standard keycode space represents ordinary HID keyboard usages and
// individual modifiers directly. The layer ranges preserve the v2 MVP
// behaviors while intentionally rejecting Vial, macro, and custom keycodes.
func decode(code uint16) (keymap.Binding, error) {
	switch {
	case code == 0:
		return keymap.Binding{Behavior: keymap.BehaviorNone}, nil
	case code >= 1 && code <= 0xdf:
		return keymap.Binding{Behavior: keymap.BehaviorKey, Param1: uint32(code)}, nil
	case code >= 0xe0 && code <= 0xe7:
		return keymap.Binding{Behavior: keymap.BehaviorModifier, Param1: 1 << (code - 0xe0)}, nil
	case code >= 0x5020 && code <= 0x50ff:
		return keymap.Binding{Behavior: keymap.BehaviorToLayer, Param1: uint32(code - 0x5020)}, nil
	case code >= 0x5220 && code <= 0x52ff:
		return keymap.Binding{Behavior: keymap.BehaviorMomentaryLayer, Param1: uint32(code - 0x5220)}, nil
	case code >= 0x5300 && code <= 0x53ff:
		return keymap.Binding{Behavior: keymap.BehaviorToggleLayer, Param1: uint32(code - 0x5300)}, nil
	default:
		return keymap.Binding{}, ErrUnsupported
	}
}

func encode(binding keymap.Binding) (uint16, error) {
	switch binding.Behavior {
	case keymap.BehaviorNone:
		return 0, nil
	case keymap.BehaviorKey:
		return uint16(binding.Param1), nil
	case keymap.BehaviorModifier:
		for bit := uint32(0); bit < 8; bit++ {
			if binding.Param1 == 1<<bit {
				return 0xe0 + uint16(bit), nil
			}
		}
	case keymap.BehaviorToLayer:
		return 0x5020 + uint16(binding.Param1), nil
	case keymap.BehaviorMomentaryLayer:
		return 0x5220 + uint16(binding.Param1), nil
	case keymap.BehaviorToggleLayer:
		return 0x5300 + uint16(binding.Param1), nil
	}
	return 0, ErrUnsupported
}
