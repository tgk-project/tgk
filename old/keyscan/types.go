package keyscan

// TGKKeyboardConfig は外部のKeyboardConfig型を表します
type TGKKeyboardConfig interface {
	GetName() string
	GetMaintainer() string
	GetVendorID() string
	GetProductID() string
	GetKeyScan() string
	GetMatrixRows() int
	GetMatrixCols() int
	GetDiodeDirection() string
	GetPushThreshold() int
	GetReleaseThreshold() int
	GetADCGain() int
	GetColChannel() []string
	GetHID() []string
	GetSplit() bool
	GetMatrixPinsRows() []string
	GetMatrixPinsCols() []string
	GetKeymap() [][]map[string]interface{}
}

// TGKKeyScanMatrix は外部のKeyScanMatrix型を表します
type TGKKeyScanMatrix interface {
	Init(config TGKKeyboardConfig) error
	Scan() bool
	GetNowPushing() [][]bool
	GetNowRelease() [][]bool
	Print() string
}
