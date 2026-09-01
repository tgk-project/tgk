package tgk

// RemapService はリマップサービスのインターフェースです
type RemapService interface {
	// Init はリマップサービスを初期化します
	Init() error
	// Remap はキーコードをリマップします
	Remap(keycode int) (int, error)
	// SetRemap はリマップ設定を追加します
	SetRemap(from, to int) error
	// ClearRemap はリマップ設定をクリアします
	ClearRemap() error
}

// NewRemapService はリマップサービスを作成します
func NewRemapService() RemapService {
	return &remapService{}
}

// remapService はリマップサービスの実装です
type remapService struct {
	remaps map[int]int
}

func (s *remapService) Init() error {
	s.remaps = make(map[int]int)
	return nil
}

func (s *remapService) Remap(keycode int) (int, error) {
	if remapped, ok := s.remaps[keycode]; ok {
		return remapped, nil
	}
	return keycode, nil
}

func (s *remapService) SetRemap(from, to int) error {
	s.remaps[from] = to
	return nil
}

func (s *remapService) ClearRemap() error {
	s.remaps = make(map[int]int)
	return nil
}
