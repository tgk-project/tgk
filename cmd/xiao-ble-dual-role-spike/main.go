//go:build xiao_ble

// xiao-ble-dual-role-spike is the Primary-side minimal reproducer for Issue
// #9. It is intentionally not a HID or split protocol implementation: it
// keeps a Host-facing GATT server connection while connecting as a GATT client
// to the companion Secondary firmware and enabling its notifications.
package main

import (
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter

	primaryServiceUUID = bluetooth.NewUUID([16]byte{0x76, 0x1d, 0x5b, 0x1c, 0x8a, 0x34, 0x45, 0x09, 0x9c, 0xa0, 0xd7, 0xd8, 0x79, 0x08, 0x00, 0x01})
	primaryNotifyUUID  = primaryServiceUUID.Replace16BitComponent(0x0802)
	primaryNotify      bluetooth.Characteristic

	secondaryServiceUUID = bluetooth.ServiceUUIDNordicUART
	secondaryRXUUID      = bluetooth.CharacteristicUUIDUARTRX
	secondaryTXUUID      = bluetooth.CharacteristicUUIDUARTTX
)

func main() {
	must("enable BLE", adapter.Enable())
	adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		if connected {
			println("connection established with", device.Address.String())
			return
		}
		println("connection disconnected")
	})

	must("add Primary service", adapter.AddService(&bluetooth.Service{
		UUID: primaryServiceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &primaryNotify,
				UUID:   primaryNotifyUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))
	advertisement := adapter.DefaultAdvertisement()
	must("configure Primary advertisement", advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "TGK-Primary-Spike",
		ServiceUUIDs: []bluetooth.UUID{primaryServiceUUID},
	}))
	must("start Primary advertisement", advertisement.Start())

	go connectSecondary()

	var sequence byte
	for {
		sequence++
		// A connected Host central that subscribes to this characteristic
		// receives these heartbeats while the Primary acts as a GATT client.
		if _, err := primaryNotify.Write([]byte{sequence}); err != nil {
			println("Host notification failed:", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func connectSecondary() {
	var secondary bluetooth.ScanResult
	println("scanning for Secondary")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() != "TGK-Secondary-Spike" || !result.HasServiceUUID(secondaryServiceUUID) {
			return
		}
		secondary = result
		must("stop scan", adapter.StopScan())
	})
	must("scan for Secondary", err)

	device, err := adapter.Connect(secondary.Address, bluetooth.ConnectionParams{})
	must("connect Secondary", err)
	services, err := device.DiscoverServices([]bluetooth.UUID{secondaryServiceUUID})
	must("discover Secondary service", err)
	if len(services) != 1 {
		panic("Secondary service discovery returned an unexpected count")
	}
	characteristics, err := services[0].DiscoverCharacteristics([]bluetooth.UUID{secondaryRXUUID, secondaryTXUUID})
	must("discover Secondary characteristics", err)
	if len(characteristics) != 2 {
		panic("Secondary characteristic discovery returned an unexpected count")
	}
	must("enable Secondary notification", characteristics[1].EnableNotifications(func(value []byte) {
		println("Secondary notification", len(value), "bytes")
	}))
	println("Secondary GATT client ready")

	var sequence byte
	for {
		sequence++
		if _, err := characteristics[0].WriteWithoutResponse([]byte{sequence}); err != nil {
			println("Secondary write failed:", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func must(action string, err error) {
	if err != nil {
		panic(action + ": " + err.Error())
	}
}
