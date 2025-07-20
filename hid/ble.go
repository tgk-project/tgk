package hid

import (
	"machine/usb/ble/hid"
)

// BLEHID はBLE HIDの実装です
type BLEHID struct {
	keyboard hid.Device
}

// NewBLEHID は新しいBLE HIDインスタンスを作成します
func NewBLEHID() *BLEHID {
	return &BLEHID{
		keyboard: hid.New(),
	}
}

// Init はBLE HIDを初期化します
func (b *BLEHID) Init() error {
	// BLE HIDの初期化処理
	return nil
}

// SendKeyboardReport はキーボードレポートを送信します
func (b *BLEHID) SendKeyboardReport(report []byte) error {
	// BLEを介してキーボードレポートを送信
	b.keyboard.SendReport(report)
	return nil
}
