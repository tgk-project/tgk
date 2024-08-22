package tgk

import (
	"context"
)

type Keyboard struct {
	HID                 []CommonHID            `json:"hid"`
	Layer               Layer                  `json:"layer"`
	KeyScan             Keyscan                `json:"keyscan"`
	MatrixPins          map[string]GpioConfig  `json:"matrix_pins"`
	KeyScanExtraConfigs map[string]interface{} `json:"keyscan_extra_configs"`
	Config              KeyboardConfig         `json:"config"`
}

func NewKeyboard(keyboardJsonData []byte) *Keyboard {

	config := NewConfigBy(keyboardJsonData)

	keyboard := Keyboard{
		Config:              config,
		MatrixPins:          ToGpios(config.MatrixPins),
		KeyScanExtraConfigs: config.KeyScanExtraConfigs,
		Layer:               NewLayer(),
	}

	// set HID objects
	for _, hid := range keyboard.Config.HID {
		switch hid {
		// case "ble":
		// 	ble := NewBLEHID()
		// 	keyboard.HID = append(keyboard.HID, ble)
		// case "usb":
		// 	usb := NewUSBHID()
		// 	keyboard.HID = append(keyboard.HID, usb)
		case "serial":
			serial := NewSerialHIDMoc()
			keyboard.HID = append(keyboard.HID, serial)
		}
	}
	return &keyboard
}

func (k *Keyboard) Start() {

	ctx := context.Background()
	print("keyboard start")

	if k.KeyScan == nil {
		rowLen := k.MatrixPins["row"].Len
		colLen := k.MatrixPins["col"].Len
		k.RegisterKeyscan(NewMXMatrix(rowLen, colLen))
	}

	k.KeyScan.Init(ctx, k)

	k.Loop(ctx)
}

func (k *Keyboard) RegisterKeyscan(obj Keyscan) {
	k.KeyScan = obj
}

func (k *Keyboard) Loop(ctx context.Context) {
	for {
		// context cancel
		select {
		case <-ctx.Done():
			return
		default:
			if k.KeyScan.Scan(ctx, k) {

				// matrix scanで何らかのキーの変化を検知した場合
				for _, event := range k.KeyScan.GetEventsAfterScan(k) {

					// key scan event
					k.SendKey(event.Row, event.Col, event.HoldFlag)
				}
			}
		}
	}
}

func (k *Keyboard) SendKey(row uint8, col uint8, holdFlag bool) {
	if holdFlag {
		k.HID[0].Down(uint16(k.GetKeyViaLayout(row, col)))
	} else {
		k.HID[0].Up(uint16(k.GetKeyViaLayout(row, col)))
	}
}

func (k *Keyboard) GetKeyViaLayout(row uint8, col uint8) uint8 {
	return KC_0
}
