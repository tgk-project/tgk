//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/tgk-project/tgk"
)

// repositorySet はリポジトリの依存関係を提供します
var repositorySet = wire.NewSet(
	tgk.NewGPIORepository,
	tgk.NewKeymapRepository,
	tgk.NewConfigRepository,
)

// serviceSet はサービスの依存関係を提供します
var serviceSet = wire.NewSet(
	tgk.NewHIDService,
	tgk.NewKeyScanService,
	tgk.NewLayerService,
	tgk.NewRemapService,
	tgk.NewSplitService,
	tgk.NewConfigService,
	tgk.NewLoggerService,
)

// InitializeTGKManager はTGKマネージャーとその依存関係を初期化します
func InitializeTGKManager() tgk.TGKManager {
	wire.Build(
		repositorySet,
		serviceSet,
		tgk.NewTGKManager,
	)
	return nil
}
