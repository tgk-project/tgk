//go:build xiao_ble

// xiao-ble-dual-role-secondary advertises the GATT client target used by the
// Issue #9 dual-role BLE spike. Flash this program to the Secondary board.
package main

import (
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter

	serviceUUID = bluetooth.ServiceUUIDNordicUART
	rxUUID      = bluetooth.CharacteristicUUIDUARTRX
	txUUID      = bluetooth.CharacteristicUUIDUARTTX
	txChar      bluetooth.Characteristic
)

func main() {
	must("enable BLE", adapter.Enable())
	must("add Secondary service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  rxUUID,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(_ bluetooth.Connection, _ int, value []byte) {
					println("secondary received", len(value), "bytes")
				},
			},
			{
				Handle: &txChar,
				UUID:   txUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))

	advertisement := adapter.DefaultAdvertisement()
	must("configure Secondary advertisement", advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "TGK-Secondary-Spike",
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	}))
	must("start Secondary advertisement", advertisement.Start())

	var sequence byte
	for {
		sequence++
		if _, err := txChar.Write([]byte{sequence}); err != nil {
			println("secondary notification failed:", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func must(action string, err error) {
	if err != nil {
		panic(action + ": " + err.Error())
	}
}
