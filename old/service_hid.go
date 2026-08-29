package tgk

import (
	"fmt"
)

// HIDService はHIDサービスのインターフェースです
type HIDService interface {
	Init(config KeyboardConfig) error
	RegisterHID(hid HIDInterface)
	SendKeyboardReport(report []byte) error
}

// HIDService はHIDサービスの実装を管理する構造体です
type hidService struct {
	hid HIDInterface
}

// NewHIDService は新しいHIDServiceを作成します
func NewHIDService() HIDService {
	return &hidService{
		hid: nil,
	}
}

// RegisterHID はHIDを登録します
func (s *hidService) RegisterHID(hid HIDInterface) {
	s.hid = hid
}

// Init は現在選択されているHIDサービスを初期化します
func (s *hidService) Init(config KeyboardConfig) error {
	if s.hid == nil {
		return fmt.Errorf("no HID implementation selected")
	}
	return s.hid.Init()
}

// SendKeyboardReport は現在選択されているHIDサービスにキーボードレポートを送信します
func (s *hidService) SendKeyboardReport(report []byte) error {
	if s.hid == nil {
		return fmt.Errorf("no HID implementation selected")
	}
	return s.hid.SendKeyboardReport(report)
}
