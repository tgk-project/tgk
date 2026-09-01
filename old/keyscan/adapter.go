package keyscan

// KeyScanMatrixAdapter はkeyscan.KeyScanMatrixをTGKKeyScanMatrixに変換するアダプターです
type KeyScanMatrixAdapter struct {
	impl KeyScanMatrix
}

// NewKeyScanMatrixAdapter は新しいKeyScanMatrixAdapterを作成します
func NewKeyScanMatrixAdapter(impl KeyScanMatrix) TGKKeyScanMatrix {
	return &KeyScanMatrixAdapter{impl: impl}
}

// Init はキースキャンマトリックスを初期化します
func (a *KeyScanMatrixAdapter) Init(config TGKKeyboardConfig) error {
	// TGKKeyboardConfigをkeyscan.KeyboardConfigに変換
	keyscanConfig := convertConfig(config)
	return a.impl.Init(keyscanConfig)
}

// Scan はマトリクススキャン関数です
func (a *KeyScanMatrixAdapter) Scan() bool {
	return a.impl.Scan()
}

// GetNowPushing は現在押されているキーの状態を返します
func (a *KeyScanMatrixAdapter) GetNowPushing() [][]bool {
	return a.impl.GetNowPushing()
}

// GetNowRelease は現在離されているキーの状態を返します
func (a *KeyScanMatrixAdapter) GetNowRelease() [][]bool {
	return a.impl.GetNowRelease()
}

// Print は現在の状態を文字列で出力します
func (a *KeyScanMatrixAdapter) Print() string {
	return a.impl.Print()
}

// convertConfig はTGKKeyboardConfigをkeyscan.KeyboardConfigに変換します
func convertConfig(config TGKKeyboardConfig) KeyboardConfig {
	return KeyboardConfig{
		Name:       config.GetName(),
		Maintainer: config.GetMaintainer(),
		VendorID:   config.GetVendorID(),
		ProductID:  config.GetProductID(),
		KeyScan:    config.GetKeyScan(),
		Matrix: Matrix{
			Rows: config.GetMatrixRows(),
			Cols: config.GetMatrixCols(),
		},
		KeyScanExtraConfigs: KeyScanExtraConfigs{
			DiodeDirection:   config.GetDiodeDirection(),
			PushThreshold:    config.GetPushThreshold(),
			ReleaseThreshold: config.GetReleaseThreshold(),
			ADCGain:          config.GetADCGain(),
			ColChannel:       config.GetColChannel(),
		},
		HID:   config.GetHID(),
		Split: config.GetSplit(),
		MatrixPins: MatrixPins{
			Rows: config.GetMatrixPinsRows(),
			Cols: config.GetMatrixPinsCols(),
		},
		Layouts: Layouts{
			Keymap: config.GetKeymap(),
		},
	}
}
