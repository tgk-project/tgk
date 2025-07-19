package tgk

// TGKManager はTGKマネージャーのインターフェースです
type TGKManager interface {
	// Init はTGKマネージャーを初期化します
	Init() error
	// Task はTGKマネージャーのタスクを実行します
	Task() error
	// SetLoadConfigName はTGKManagerが読み込むファイル名を指定します
	SetLoadConfigName(configFileName string) error
	// SetLogLevel はログレベルを設定します
	SetLogLevel(level LogLevel)
	// GetLoggerService はロガーサービスを取得します
	GetLoggerService() LoggerService
}

// tgkManager はTGKマネージャーの実装です
type tgkManager struct {
	hidService     HIDService
	keyScanService KeyScanService
	layerService   LayerService
	remapService   RemapService
	splitService   SplitService
	configService  ConfigService
	loggerService  LoggerService
}

// NewTGKManager はTGKマネージャーを作成します
func NewTGKManager(hidService HIDService, keyScanService KeyScanService,
	layerService LayerService, remapService RemapService, splitService SplitService, configService ConfigService, loggerService LoggerService) TGKManager {
	return &tgkManager{
		hidService:     hidService,
		keyScanService: keyScanService,
		layerService:   layerService,
		remapService:   remapService,
		splitService:   splitService,
		configService:  configService,
		loggerService:  loggerService,
	}
}

func (m *tgkManager) Init() error {
	m.loggerService.Debug("TGKManager初期化開始")

	m.loggerService.Debug("設定サービス初期化中...")
	if err := m.configService.Init(); err != nil {
		m.loggerService.Error("設定サービス初期化エラー: %v", err)
		return err
	}
	m.loggerService.Debug("設定サービス初期化完了")

	config := m.configService.GetConfig()
	m.loggerService.Debug("設定読み込み完了: Matrix=%dx%d", config.Matrix.Rows, config.Matrix.Cols)

	// ログレベルの設定
	if config.LogLevel != "" {
		logLevel := ParseLogLevel(config.LogLevel)
		m.loggerService.SetLevel(logLevel)
		m.loggerService.Info("ログレベルを設定しました: %s", logLevel.String())
	}

	m.loggerService.Debug("HIDサービス初期化中...")
	if err := m.hidService.Init(config); err != nil {
		m.loggerService.Error("HIDサービス初期化エラー: %v", err)
		return err
	}
	m.loggerService.Debug("HIDサービス初期化完了")

	m.loggerService.Debug("キースキャンサービス初期化中...")
	if err := m.keyScanService.Init(config); err != nil {
		m.loggerService.Error("キースキャンサービス初期化エラー: %v", err)
		return err
	}
	m.loggerService.Debug("キースキャンサービス初期化完了")

	m.loggerService.Debug("レイヤーサービス初期化中...")
	if err := m.layerService.Init(config); err != nil {
		m.loggerService.Error("レイヤーサービス初期化エラー: %v", err)
		return err
	}
	m.loggerService.Debug("レイヤーサービス初期化完了")

	// if err := m.remapService.Init(config); err != nil {
	// 	return err
	// }
	// if err := m.splitService.Init(config); err != nil {
	// 	return err
	// }

	m.loggerService.Info("TGKManager初期化完了")
	return nil
}

func (m *tgkManager) SetLoadConfigName(configFileName string) error {
	return m.configService.SetLoadConfigName(configFileName)
}

func (m *tgkManager) SetLogLevel(level LogLevel) {
	m.loggerService.SetLevel(level)
}

func (m *tgkManager) GetLoggerService() LoggerService {
	return m.loggerService
}

func (m *tgkManager) Task() error {
	m.loggerService.Debug("タスク実行開始")

	// キースキャン
	isMatrixUpdate := m.keyScanService.Scan()
	m.loggerService.Debug("キースキャン完了: マトリックス更新=%t", isMatrixUpdate)

	// キーマップ処理
	if isMatrixUpdate {
		m.loggerService.Debug("マトリックス更新検出 - キーマップ処理開始")
		// TODO: 実際のキーマップ処理を実装
		// 現在は基本的なループ構造のみ
		config := m.configService.GetConfig()
		for row := 0; row < config.Matrix.Rows; row++ {
			for col := 0; col < config.Matrix.Cols; col++ {
				keycode, err := m.layerService.GetKeymap(row, col)
				if err != nil {
					m.loggerService.Error("キーマップ取得エラー [%d,%d]: %v", row, col, err)
					return err
				}
				// キーコードの処理（現在は未使用変数を回避するため_で受ける）
				if keycode != 0 {
					m.loggerService.Debug("キー検出 [%d,%d]: keycode=%d", row, col, keycode)
				}
				_ = keycode
			}
		}
		m.loggerService.Debug("キーマップ処理完了")
	}

	return nil
}
