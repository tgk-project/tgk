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

// KeyScanMatrixFactory はキースキャンマトリックスのファクトリー関数の型です
type KeyScanMatrixFactory func() KeyScanMatrix

// registeredFactories は登録されたファクトリー関数のマップです
var registeredFactories = map[string]KeyScanMatrixFactory{
	"mx": func() KeyScanMatrix { return NewMXKeyScan() },
	"ec": func() KeyScanMatrix { return NewECKeyScan() },
}

// RegisterKeyScanMatrix は新しいキースキャンマトリックスタイプを登録します
func RegisterKeyScanMatrix(typeName string, factory KeyScanMatrixFactory) {
	registeredFactories[typeName] = factory
}

// ErrUnknownKeyScanType は未知のキースキャンタイプが指定された場合のエラーです
type ErrUnknownKeyScanType struct {
	Type string
}

func (e ErrUnknownKeyScanType) Error() string {
	return "unknown keyscan type: " + e.Type
}

// NewKeyScanMatrix はキースキャンタイプに基づいて適切なKeyScanMatrixを作成します
// 未知のタイプが指定された場合はデフォルト（MX）を使用します
func NewKeyScanMatrix(keyscanType string) KeyScanMatrix {
	factory, exists := registeredFactories[keyscanType]
	if !exists {
		// デフォルトはMX
		factory = func() KeyScanMatrix { return NewMXKeyScan() }
	}

	return factory()
}

// NewKeyScanMatrixWithError はキースキャンタイプに基づいて適切なKeyScanMatrixを作成します
// 未知のタイプが指定された場合はエラーを返します
func NewKeyScanMatrixWithError(keyscanType string) (KeyScanMatrix, error) {
	factory, exists := registeredFactories[keyscanType]
	if !exists {
		return nil, ErrUnknownKeyScanType{Type: keyscanType}
	}

	return factory(), nil
}

// NewTGKKeyScanMatrix はキースキャンタイプに基づいて適切なTGKKeyScanMatrixを作成します
func NewTGKKeyScanMatrix(keyscanType string) TGKKeyScanMatrix {
	impl := NewKeyScanMatrix(keyscanType)
	return NewKeyScanMatrixAdapter(impl)
}
