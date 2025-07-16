package tgk

import "github.com/Diwamoto/tgk/keyscan"

// KeyScanService はキースキャンサービスのインターフェースです
type KeyScanService interface {
	// Init はキースキャンサービスを初期化します
	Init(config KeyboardConfig) error
	// Scan はキーマトリックスをスキャンし、押されているキーの情報を返します
	Scan() bool
	// Print は現在の状態を文字列で出力します
	Print() string
}

type keyScanService struct {
	keyScan    KeyScanMatrix
	nowPushing [][]bool
	nowRelease [][]bool
}

func NewKeyScanService() KeyScanService {
	return &keyScanService{}
}

func (s *keyScanService) Init(config KeyboardConfig) error {
	// configからkeyscanタイプを取得して適切なKeyScanMatrixを作成
	s.keyScan = keyscan.NewKeyScanMatrix(config.KeyScan)

	// tgk.KeyboardConfigをkeyscan.KeyboardConfigに変換
	keyscanConfig := keyscan.KeyboardConfig{
		Name:       config.Name,
		Maintainer: config.Maintainer,
		VendorID:   config.VendorID,
		ProductID:  config.ProductID,
		KeyScan:    config.KeyScan,
		Matrix: keyscan.Matrix{
			Rows: config.Matrix.Rows,
			Cols: config.Matrix.Cols,
		},
		KeyScanExtraConfigs: keyscan.KeyScanExtraConfigs{
			DiodeDirection:   config.KeyScanExtraConfigs.DiodeDirection,
			PushThreshold:    config.KeyScanExtraConfigs.PushThreshold,
			ReleaseThreshold: config.KeyScanExtraConfigs.ReleaseThreshold,
			ADCGain:          config.KeyScanExtraConfigs.ADCGain,
			ColChannel:       config.KeyScanExtraConfigs.ColChannel,
		},
		HID:   config.HID,
		Split: config.Split,
		MatrixPins: keyscan.MatrixPins{
			Rows: config.MatrixPins.Rows,
			Cols: config.MatrixPins.Cols,
		},
		Layouts: keyscan.Layouts{
			Keymap: config.Layouts.Keymap,
		},
	}

	return s.keyScan.Init(keyscanConfig)
}

func (s *keyScanService) Scan() bool {
	isMatrixUpdate := s.keyScan.Scan()

	if isMatrixUpdate {
		s.nowPushing = s.keyScan.GetNowPushing()
		s.nowRelease = s.keyScan.GetNowRelease()
	}

	return isMatrixUpdate
}

func (s *keyScanService) Print() string {
	return s.keyScan.Print()
}
