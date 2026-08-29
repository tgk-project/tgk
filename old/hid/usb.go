//go:build tinygo
// +build tinygo

package hid

// TinyGoでは machine/usb/hid/keyboard パッケージを使用
// 通常のGoコンパイラでは利用できないため、ビルド制約を使用

import (
	kbd "machine/usb/hid/keyboard"
)

// machine/usb/hid/keyboard のラッパー
type USBHID struct {
	keyboard kbd.Device
}

func NewUSBHID() *USBHID {
	return &USBHID{
		keyboard: kbd.Port(),
	}
}

func (s *USBHID) Init() error {
	return nil
}

func (s *USBHID) SendKeyboardReport(report []byte) error {
	// キーボードレポートの送信
	s.keyboard.SendReport(report)
	return nil
}

func (s *USBHID) SendMouseReport(report []byte) error {
	// マウスレポートの送信
	// TODO: 実装
	return nil
}
