package tgk

// KeyboardConfigのメソッド実装
func (c KeyboardConfig) GetName() string {
	return c.Name
}

func (c KeyboardConfig) GetMaintainer() string {
	return c.Maintainer
}

func (c KeyboardConfig) GetVendorID() string {
	return c.VendorID
}

func (c KeyboardConfig) GetProductID() string {
	return c.ProductID
}

func (c KeyboardConfig) GetKeyScan() string {
	return c.KeyScan
}

func (c KeyboardConfig) GetMatrixRows() int {
	return c.Matrix.Rows
}

func (c KeyboardConfig) GetMatrixCols() int {
	return c.Matrix.Cols
}

func (c KeyboardConfig) GetDiodeDirection() string {
	return c.KeyScanExtraConfigs.DiodeDirection
}

func (c KeyboardConfig) GetPushThreshold() int {
	return c.KeyScanExtraConfigs.PushThreshold
}

func (c KeyboardConfig) GetReleaseThreshold() int {
	return c.KeyScanExtraConfigs.ReleaseThreshold
}

func (c KeyboardConfig) GetADCGain() int {
	return c.KeyScanExtraConfigs.ADCGain
}

func (c KeyboardConfig) GetColChannel() []string {
	return c.KeyScanExtraConfigs.ColChannel
}

func (c KeyboardConfig) GetHID() []string {
	return c.HID
}

func (c KeyboardConfig) GetSplit() bool {
	return c.Split
}

func (c KeyboardConfig) GetMatrixPinsRows() []string {
	return c.MatrixPins.Rows
}

func (c KeyboardConfig) GetMatrixPinsCols() []string {
	return c.MatrixPins.Cols
}

func (c KeyboardConfig) GetKeymap() [][]map[string]interface{} {
	return c.Layouts.Keymap
}
