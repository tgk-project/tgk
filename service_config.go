package tgk

// ConfigService は設定サービスのインターフェースです
type ConfigService interface {
	// Init は設定サービスを初期化します
	Init() error
	// GetConfig は設定を取得します
	GetConfig() KeyboardConfig
	// SetLoadConfigName はTGKManagerが読み込むファイル名を指定します
	SetLoadConfigName(configFileName string) error
	// Load は設定を読み込みます
	Load() error
}

// configService は設定サービスの実装です
type configService struct {
	configRepository ConfigRepository
	loadConfigName   string
}

// NewConfigService は設定サービスを作成します
func NewConfigService(configRepository ConfigRepository) ConfigService {
	return &configService{
		configRepository: configRepository,
		loadConfigName:   "keyboard.json", // デフォルトはkeyboard.json
	}
}

func (s *configService) Init() error {
	if err := s.Load(); err != nil {
		return err
	}
	return nil
}

func (s *configService) GetConfig() KeyboardConfig {
	return s.configRepository.GetConfig()
}

func (s *configService) SetLoadConfigName(configFileName string) error {
	s.loadConfigName = configFileName
	return nil
}

func (s *configService) Load() error {
	return s.configRepository.Load(s.loadConfigName)
}
