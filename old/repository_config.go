package tgk

import (
	"encoding/json"
	"os"
)

// ConfigRepository は設定リポジトリのインターフェースです
type ConfigRepository interface {
	// Init は設定リポジトリを初期化します
	Init() error
	// Load は設定を読み込みます
	Load(configFileName string) error
	// GetConfig は設定を返します
	GetConfig() KeyboardConfig
}

// configRepository は設定リポジトリの実装です
type configRepository struct {
	config KeyboardConfig
}

// NewConfigRepository は設定リポジトリを作成します
func NewConfigRepository() ConfigRepository {
	return &configRepository{}
}

func (r *configRepository) Init() error {
	r.config = KeyboardConfig{}
	return nil
}

func (r *configRepository) Load(configFileName string) error {
	// keyboard.jsonの読み込み
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	// JSONのデコード
	if err := json.Unmarshal(data, &r.config); err != nil {
		return err
	}

	return nil
}

func (r *configRepository) GetConfig() KeyboardConfig {
	return r.config
}
