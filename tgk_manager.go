package tgk

// TGKManager はTGKマネージャーのインターフェースです
type TGKManager interface {
	// Init はTGKマネージャーを初期化します
	Init() error
	// Task はTGKマネージャーのタスクを実行します
	Task() error
	// SetLoadConfigName はTGKManagerが読み込むファイル名を指定します
	SetLoadConfigName(configFileName string) error
}

// tgkManager はTGKマネージャーの実装です
type tgkManager struct {
	hidService     HIDService
	keyScanService KeyScanService
	layerService   LayerService
	remapService   RemapService
	splitService   SplitService
	configService  ConfigService
}

// NewTGKManager はTGKマネージャーを作成します
func NewTGKManager(hidService HIDService, keyScanService KeyScanService,
	layerService LayerService, remapService RemapService, splitService SplitService, configService ConfigService) TGKManager {
	return &tgkManager{
		hidService:     hidService,
		keyScanService: keyScanService,
		layerService:   layerService,
		remapService:   remapService,
		splitService:   splitService,
		configService:  configService,
	}
}

func (m *tgkManager) Init() error {

	if err := m.configService.Init(); err != nil {
		return err
	}

	config := m.configService.GetConfig()

	if err := m.hidService.Init(config); err != nil {
		return err
	}
	if err := m.keyScanService.Init(config); err != nil {
		return err
	}
	if err := m.layerService.Init(config); err != nil {
		return err
	}
	// if err := m.remapService.Init(config); err != nil {
	// 	return err
	// }
	// if err := m.splitService.Init(config); err != nil {
	// 	return err
	// }
	return nil
}

func (m *tgkManager) SetLoadConfigName(configFileName string) error {
	return m.configService.SetLoadConfigName(configFileName)
}

func (m *tgkManager) Task() error {
	// キースキャン
	isMatrixUpdate := m.keyScanService.Scan()

	// キーマップ処理
	if isMatrixUpdate {
		// TODO: 実際のキーマップ処理を実装
		// 現在は基本的なループ構造のみ
		config := m.configService.GetConfig()
		for row := 0; row < config.Matrix.Rows; row++ {
			for col := 0; col < config.Matrix.Cols; col++ {
				keycode, err := m.layerService.GetKeymap(row, col)
				if err != nil {
					return err
				}
				// キーコードの処理（現在は未使用変数を回避するため_で受ける）
				_ = keycode
			}
		}
	}

	return nil
}
