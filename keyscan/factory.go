package keyscan

// KeyScanMatrix はキースキャンマトリックスのインターフェースです
type KeyScanMatrix interface {
	// Init はキースキャンマトリックスを初期化します
	Init(config KeyboardConfig) error
	// Scan はマトリクススキャン関数です。
	Scan() bool
	// GetNowPushing は現在押されているキーの状態を返します
	GetNowPushing() [][]bool
	// GetNowRelease は現在離されているキーの状態を返します
	GetNowRelease() [][]bool
	// Print は現在の状態を文字列で出力します
	Print() string
}

// KeyboardConfig はキーボードの設定を表します（循環インポート回避のため再定義）
type KeyboardConfig struct {
	Name                string              `json:"name"`
	Maintainer          string              `json:"maintainer"`
	VendorID            string              `json:"vendorId"`
	ProductID           string              `json:"productId"`
	KeyScan             string              `json:"keyscan"`
	Matrix              Matrix              `json:"matrix"`
	KeyScanExtraConfigs KeyScanExtraConfigs `json:"keyscan_extra_configs"`
	HID                 []string            `json:"hid"`
	Split               bool                `json:"split"`
	MatrixPins          MatrixPins          `json:"matrix_pins"`
	Layouts             Layouts             `json:"layouts"`
}

// Matrix はマトリックスの設定を表します
type Matrix struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

// KeyScanExtraConfigs はキースキャンの追加設定を表します
type KeyScanExtraConfigs struct {
	DiodeDirection   string   `json:"diode_direction"`
	PushThreshold    int      `json:"push_threshold"`
	ReleaseThreshold int      `json:"release_threshold"`
	ADCGain          int      `json:"adc_gain"`
	ColChannel       []string `json:"col_channel"`
}

// MatrixPins はマトリックスのピン設定を表します
type MatrixPins struct {
	Rows []string `json:"rows"`
	Cols []string `json:"cols"`
}

// Layouts はレイアウトの設定を表します
type Layouts struct {
	Keymap [][]map[string]interface{} `json:"keymap"`
}

// NewKeyScanMatrix はキースキャンマトリックスを作成します
func NewKeyScanMatrix(keyscanType string) KeyScanMatrix {
	switch keyscanType {
	case "mx":
		return NewMXKeyScan()
	case "ec":
		return NewECKeyScan()
	default:
		// デフォルトはMX
		return NewMXKeyScan()
	}
}
