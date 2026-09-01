//go:build xiao_ble

// xiao-ble-usb is a firmware build probe for the XIAO BLE USB composite HID
// adapter. It intentionally contains no VIA command implementation.
package main

import "github.com/tgk-project/tgk/platform/xiao_ble"

func main() {
	_, _ = xiaoble.NewUSBHID(4)
	for {
	}
}
