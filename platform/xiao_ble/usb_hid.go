//go:build xiao_ble

package xiaoble

import (
	"errors"

	"machine"
	"machine/usb"
	"machine/usb/descriptor"

	usbtransport "github.com/tgk-project/tgk/keyboard/transport/usb"
)

const (
	keyboardEndpoint = usb.HID_ENDPOINT_IN  // endpoint 4
	consumerEndpoint = usb.HID_ENDPOINT_OUT // endpoint 5, used as IN here
	rawInEndpoint    = 6
	rawOutEndpoint   = 7
)

var ErrUSBHIDConfigured = errors.New("XIAO BLE USB HID is already configured")

// usbHID owns the TinyGo callback lifetime. TinyGo's ConfigureUSBEndpoint API
// requires package-level function callbacks, so the configured transport is
// retained here instead of exposing a local adapter implementation to callers.
var usbHID *usbtransport.Device

// NewUSBHID configures a composite USB device with independent Boot Keyboard,
// Consumer Control, and vendor-defined RAW HID interfaces. The returned
// transport is suitable for keyboard/engine.HostTransport; RAW HID data is
// intentionally available only as opaque 32-byte packets for Issue #11.
func NewUSBHID(rawQueueCapacity int) (*usbtransport.Device, error) {
	if usbHID != nil {
		return nil, ErrUSBHIDConfigured
	}
	device, err := usbtransport.New(tinyGoEndpoint{}, rawQueueCapacity)
	if err != nil {
		return nil, err
	}
	usbHID = device
	machine.ConfigureUSBEndpoint(compositeDescriptor,
		[]usb.EndpointConfig{
			{Index: keyboardEndpoint, IsIn: true, Type: usb.ENDPOINT_TYPE_INTERRUPT},
			{Index: consumerEndpoint, IsIn: true, Type: usb.ENDPOINT_TYPE_INTERRUPT},
			{Index: rawInEndpoint, IsIn: true, Type: usb.ENDPOINT_TYPE_INTERRUPT},
			{Index: rawOutEndpoint, IsIn: false, Type: usb.ENDPOINT_TYPE_INTERRUPT, RxHandler: receiveRaw},
		},
		[]usb.SetupConfig{
			{Index: 0, Handler: hidSetup},
			{Index: 1, Handler: hidSetup},
			{Index: 2, Handler: hidSetup},
		},
	)
	return device, nil
}

type tinyGoEndpoint struct{}

func (tinyGoEndpoint) Ready() bool { return machine.USBDev.InitEndpointComplete }

func (tinyGoEndpoint) SendKeyboard(packet []byte) bool {
	return machine.SendUSBInPacket(keyboardEndpoint, packet)
}

func (tinyGoEndpoint) SendConsumer(packet []byte) bool {
	return machine.SendUSBInPacket(consumerEndpoint, packet)
}

func (tinyGoEndpoint) SendRaw(packet []byte) bool {
	return machine.SendUSBInPacket(rawInEndpoint, packet)
}

func receiveRaw(packet []byte) {
	if usbHID != nil {
		usbHID.HandleRawOutput(packet)
	}
}

// hidSetup acknowledges the class requests needed for the three static HID
// interfaces. Report descriptors are served by TinyGo from compositeDescriptor.
func hidSetup(setup usb.Setup) bool {
	if setup.BmRequestType == usb.SET_REPORT_TYPE &&
		(setup.BRequest == usb.SET_IDLE || setup.BRequest == usb.SET_PROTOCOL) {
		machine.SendZlp()
		return true
	}
	return false
}

var compositeDescriptor = descriptor.Descriptor{
	Device: []byte{
		18, descriptor.TypeDevice, 0x00, 0x02, 0x00, 0x00, 0x00, 64,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 1, 2, 0, 1,
	},
	Configuration: []byte{
		9, descriptor.TypeConfiguration, 91, 0, 3, 1, 0, usb.CONFIG_BUS_POWERED, 50,

		// Interface 0: Boot Keyboard on endpoint 4 IN.
		9, descriptor.TypeInterface, 0, 0, 1, usb.DEVICE_CLASS_HUMAN_INTERFACE, 1, 1, 0,
		9, descriptor.TypeClassHID, 0x11, 0x01, 0, 1, descriptor.TypeHIDReport, 63, 0,
		7, descriptor.TypeEndpoint, usb.EndpointIn | keyboardEndpoint, usb.ENDPOINT_TYPE_INTERRUPT, 8, 0, 10,

		// Interface 1: Consumer Control on endpoint 5 IN.
		9, descriptor.TypeInterface, 1, 0, 1, usb.DEVICE_CLASS_HUMAN_INTERFACE, 0, 0, 0,
		9, descriptor.TypeClassHID, 0x11, 0x01, 0, 1, descriptor.TypeHIDReport, 23, 0,
		7, descriptor.TypeEndpoint, usb.EndpointIn | consumerEndpoint, usb.ENDPOINT_TYPE_INTERRUPT, 2, 0, 10,

		// Interface 2: vendor-defined RAW HID on endpoints 6 IN and 7 OUT.
		9, descriptor.TypeInterface, 2, 0, 2, usb.DEVICE_CLASS_HUMAN_INTERFACE, 0, 0, 0,
		9, descriptor.TypeClassHID, 0x11, 0x01, 0, 1, descriptor.TypeHIDReport, 27, 0,
		7, descriptor.TypeEndpoint, usb.EndpointIn | rawInEndpoint, usb.ENDPOINT_TYPE_INTERRUPT, usbtransport.RawReportSize, 0, 1,
		7, descriptor.TypeEndpoint, usb.EndpointOut | rawOutEndpoint, usb.ENDPOINT_TYPE_INTERRUPT, usbtransport.RawReportSize, 0, 1,
	},
	HID: map[uint16][]byte{
		0: {
			0x05, 0x01, 0x09, 0x06, 0xa1, 0x01,
			0x05, 0x07, 0x19, 0xe0, 0x29, 0xe7, 0x15, 0x00, 0x25, 0x01, 0x75, 0x01, 0x95, 0x08, 0x81, 0x02,
			0x95, 0x01, 0x75, 0x08, 0x81, 0x03,
			0x95, 0x06, 0x75, 0x08, 0x15, 0x00, 0x25, 0x65, 0x05, 0x07, 0x19, 0x00, 0x29, 0x65, 0x81, 0x00,
			0xc0,
		},
		1: {
			0x05, 0x0c, 0x09, 0x01, 0xa1, 0x01,
			0x15, 0x00, 0x26, 0xff, 0x03, 0x19, 0x00, 0x2a, 0xff, 0x03, 0x75, 0x10, 0x95, 0x01, 0x81, 0x00,
			0xc0,
		},
		2: {
			0x06, 0x60, 0xff, 0x09, 0x61, 0xa1, 0x01,
			0x09, 0x62, 0x15, 0x00, 0x26, 0xff, 0x00, 0x75, 0x08, 0x95, usbtransport.RawReportSize, 0x81, 0x02,
			0x09, 0x63, 0x95, usbtransport.RawReportSize, 0x91, 0x02,
			0xc0,
		},
	},
}
