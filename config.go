package tgk

import (
	_ "embed"

	"encoding/json"
)

type KeyboardConfig struct {
	Name                string                 `json:"name"`
	Maintainer          string                 `json:"maintainer"`
	VendorID            string                 `json:"vendorId"`
	ProductID           string                 `json:"productId"`
	KeyScan             string                 `json:"keyscan"`
	Matrix              map[string]interface{} `json:"matrix"`
	KeyScanExtraConfigs map[string]interface{} `json:"keyscan_extra_configs"`
	HID                 []string               `json:"hid"`
	Split               bool                   `json:"split"`
	MatrixPins          map[string][]string    `json:"matrix_pins"`
	Layouts             map[string]interface{} `json:"layouts"`
}

func NewConfigBy(keyboardJsonData []byte) KeyboardConfig {

	config := KeyboardConfig{}
	err := json.Unmarshal(keyboardJsonData, &config)
	if err != nil {
		panic(err)
	}

	return config
}
