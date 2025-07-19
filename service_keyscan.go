package tgk

import "github.com/tgk-project/tgk/keyscan"

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
	keyScan    keyscan.TGKKeyScanMatrix
	nowPushing [][]bool
	nowRelease [][]bool
}

func NewKeyScanService() KeyScanService {
	return &keyScanService{}
}

func (s *keyScanService) Init(config KeyboardConfig) error {
	// configからkeyscanタイプを取得して適切なKeyScanMatrixを作成
	// 新しいアダプターを使用して、tgk.KeyScanMatrixを直接取得
	s.keyScan = keyscan.NewTGKKeyScanMatrix(config.KeyScan)

	// 直接tgk.KeyboardConfigを渡せるようになった
	return s.keyScan.Init(config)
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
