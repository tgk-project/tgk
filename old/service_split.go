package tgk

// SplitService はスプリットサービスのインターフェースです
type SplitService interface {
	// Init はスプリットサービスを初期化します
	Init() error
	// IsMaster はマスター側かどうかを返します
	IsMaster() bool
	// Send はデータを送信します
	Send(data []byte) error
	// Receive はデータを受信します
	Receive() ([]byte, error)
}

// NewSplitService はスプリットサービスを作成します
func NewSplitService() SplitService {
	return &splitService{}
}

// splitService はスプリットサービスの実装です
type splitService struct {
	isMaster bool
}

func (s *splitService) Init() error {
	// TODO: 実装
	return nil
}

func (s *splitService) IsMaster() bool {
	return s.isMaster
}

func (s *splitService) Send(data []byte) error {
	// TODO: 実装
	return nil
}

func (s *splitService) Receive() ([]byte, error) {
	// TODO: 実装
	return nil, nil
}
