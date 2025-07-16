package tgk

// KeymapRepository はキーマップリポジトリのインターフェースです
type KeymapRepository interface {
	// Init はキーマップリポジトリを初期化します
	Init() error
	// Load はキーマップを読み込みます
	Load() error
	// Save はキーマップを保存します
	Save() error
	// GetKeymap はキーマップを返します
	GetKeymap() [][]int
	// SetKeymap はキーマップを設定します
	SetKeymap(keymap [][]int) error
}

// KeymapConfig はキーマップの設定を表します
type KeymapConfig struct {
	Layers [][]int
}

// keymapRepository はキーマップリポジトリの実装です
type keymapRepository struct {
	config     *KeymapConfig
	configRepo ConfigRepository
}

// NewKeymapRepository はキーマップリポジトリを作成します
func NewKeymapRepository(configRepo ConfigRepository) KeymapRepository {
	return &keymapRepository{
		config:     &KeymapConfig{},
		configRepo: configRepo,
	}
}

func (r *keymapRepository) Init() error {
	return nil
}

func (r *keymapRepository) Load() error {
	// // 設定の読み込み
	// config := r.configRepo.GetConfig()

	// // レイアウトの取得
	// layouts, ok := config["layouts"].(map[string]interface{})
	// if !ok {
	// 	return fmt.Errorf("invalid layouts format")
	// }

	// // キーマップの取得
	// keymap, ok := layouts["keymap"].([][]map[string]interface{})
	// if !ok {
	// 	return fmt.Errorf("invalid keymap format")
	// }

	// // キーマップの変換
	// matrix, ok := config["matrix"].(map[string]interface{})
	// if !ok {
	// 	return fmt.Errorf("invalid matrix format")
	// }

	// rows := int(matrix["rows"].(float64))
	// cols := int(matrix["cols"].(float64))

	// // レイヤー数はデフォルトで1
	// r.config.Layers = make([][]int, 1)
	// r.config.Layers[0] = make([]int, rows*cols)

	// // キーマップの解析
	// for i, row := range keymap {
	// 	for _, key := range row {
	// 		// キーの位置情報を取得
	// 		x, ok := key["x"].(float64)
	// 		if !ok {
	// 			x = 0
	// 		}
	// 		y, ok := key["y"].(float64)
	// 		if !ok {
	// 			y = float64(i)
	// 		}

	// 		// キーコードを取得
	// 		code, ok := key["code"].(float64)
	// 		if !ok {
	// 			code = 0
	// 		}

	// 		// キーマップに設定
	// 		pos := int(y)*cols + int(x)
	// 		if pos >= 0 && pos < len(r.config.Layers[0]) {
	// 			r.config.Layers[0][pos] = int(code)
	// 		}
	// 	}
	// }

	return nil
}

func (r *keymapRepository) Save() error {
	// // 設定の取得
	// config := r.configRepo.GetConfig()

	// // レイアウトの取得
	// layouts, ok := config["layouts"].(map[string]interface{})
	// if !ok {
	// 	layouts = make(map[string]interface{})
	// 	config["layouts"] = layouts
	// }

	// // キーマップの変換
	// keymap := make([][]map[string]interface{}, len(r.config.Layers))
	// for i, layer := range r.config.Layers {
	// 	keymap[i] = make([]map[string]interface{}, len(layer))
	// 	for j, code := range layer {
	// 		keymap[i][j] = map[string]interface{}{
	// 			"code": code,
	// 			"x":    j % len(layer),
	// 			"y":    j / len(layer),
	// 		}
	// 	}
	// }

	// // キーマップの設定
	// layouts["keymap"] = keymap

	// // 設定の保存
	// return r.configRepo.SetConfig(config)
	return nil
}

func (r *keymapRepository) GetKeymap() [][]int {
	return r.config.Layers
}

func (r *keymapRepository) SetKeymap(keymap [][]int) error {
	r.config.Layers = keymap
	return nil
}
