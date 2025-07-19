# キースキャンタイプの追加ガイド

このドキュメントでは、TGKプロジェクトに新しいキースキャンタイプを追加する方法について説明します。

## 概要

TGKプロジェクトでは、キーボードの設定ファイル（keyboard.json）で指定されたキースキャンタイプに基づいて、適切なキースキャンマトリックスの実装を動的に読み込む機能を提供しています。現在、以下のキースキャンタイプがサポートされています：

- `mx`: 一般的なメカニカルスイッチ用のマトリックススキャン
- `ec`: 静電容量式スイッチ用のマトリックススキャン

新しいキースキャンタイプを追加するには、以下の手順に従ってください。

## 手順

### 1. 新しいキースキャンマトリックスの実装

まず、`keyscan` パッケージに新しいキースキャンマトリックスの実装を作成します。例えば、`keyscan/new_type.go` というファイルを作成します：

```go
package keyscan

import (
	"fmt"
	"machine"
)

// KeyScanNewType は新しいタイプのキースキャンマトリックスの実装です
type KeyScanNewType struct {
	rowPins    []machine.Pin
	colPins    []machine.Pin
	nowPushing [][]bool
	nowRelease [][]bool
	config     KeyboardConfig
	// 必要に応じて追加のフィールドを定義
}

// NewNewTypeKeyScan は新しいタイプのキースキャンマトリックスを作成します
func NewNewTypeKeyScan() *KeyScanNewType {
	return &KeyScanNewType{
		// デフォルト値を設定
		rowPins: []machine.Pin{machine.D12, machine.D11, machine.D10, machine.D9, machine.D6},
		colPins: []machine.Pin{machine.D5, machine.D4, machine.D3, machine.D2, machine.D1, machine.D0},
	}
}

// Init はキースキャンマトリックスを初期化します
func (s *KeyScanNewType) Init(config KeyboardConfig) error {
	s.config = config

	// マトリックスの状態を初期化
	s.nowPushing = make([][]bool, config.Matrix.Rows)
	s.nowRelease = make([][]bool, config.Matrix.Rows)
	for i := range s.nowPushing {
		s.nowPushing[i] = make([]bool, config.Matrix.Cols)
		s.nowRelease[i] = make([]bool, config.Matrix.Cols)
	}

	// ピン設定を動的に読み込む
	if len(config.MatrixPins.Rows) > 0 && len(config.MatrixPins.Cols) > 0 {
		// 設定ファイルからピン設定を読み込む
		s.rowPins = make([]machine.Pin, len(config.MatrixPins.Rows))
		s.colPins = make([]machine.Pin, len(config.MatrixPins.Cols))
		
		for i, pinName := range config.MatrixPins.Rows {
			pin, err := stringToPinWithError(pinName)
			if err != nil {
				return fmt.Errorf("invalid row pin: %w", err)
			}
			s.rowPins[i] = pin
		}
		
		for i, pinName := range config.MatrixPins.Cols {
			pin, err := stringToPinWithError(pinName)
			if err != nil {
				return fmt.Errorf("invalid column pin: %w", err)
			}
			s.colPins[i] = pin
		}
	}

	// ピンの初期化
	// 必要に応じてピンの設定を行う

	return nil
}

// Scan はマトリクススキャン関数です
func (s *KeyScanNewType) Scan() bool {
	var isMatrixUpdate bool

	// 実際のキーマトリックススキャンを実装
	// 必要に応じてハードウェア固有のロジックを実装

	return isMatrixUpdate
}

// GetNowPushing は現在押されているキーの状態を返します
func (s *KeyScanNewType) GetNowPushing() [][]bool {
	return s.nowPushing
}

// GetNowRelease は現在離されているキーの状態を返します
func (s *KeyScanNewType) GetNowRelease() [][]bool {
	return s.nowRelease
}

// Print は現在の状態を文字列で出力します
func (s *KeyScanNewType) Print() string {
	paper := ""
	for row := 0; row < len(s.rowPins); row++ {
		for col := 0; col < len(s.colPins); col++ {
			if s.nowPushing[row][col] {
				paper += "1 "
			} else {
				paper += "0 "
			}
		}
		paper += "\n"
	}
	return paper
}
```

### 2. ファクトリー関数の登録

次に、`keyscan/factory.go` ファイルに新しいキースキャンタイプのファクトリー関数を登録します。`init` 関数を追加して、アプリケーションの起動時に自動的に登録されるようにします：

```go
func init() {
	// 新しいキースキャンタイプを登録
	RegisterKeyScanMatrix("new_type", func() KeyScanMatrix { return NewNewTypeKeyScan() })
}
```

または、アプリケーションの初期化時に明示的に登録することもできます：

```go
// アプリケーションの初期化時
keyscan.RegisterKeyScanMatrix("new_type", func() keyscan.KeyScanMatrix { return keyscan.NewNewTypeKeyScan() })
```

### 3. テストの作成

新しいキースキャンタイプのテストを作成します。`tests/keyscan_new_type_test.go` というファイルを作成します：

```go
package tests

import (
	"testing"

	"github.com/tgk-project/tgk"
	"github.com/tgk-project/tgk/keyscan"
)

func TestNewTypeKeyScan(t *testing.T) {
	// 新しいキースキャンタイプのインスタンスを作成
	matrix := keyscan.NewNewTypeKeyScan()
	
	// テスト用の設定を作成
	config := keyscan.KeyboardConfig{
		Name:     "TestKeyboard",
		KeyScan:  "new_type",
		Matrix:   keyscan.Matrix{Rows: 2, Cols: 2},
		HID:      []string{"usb"},
		Split:    false,
		KeyScanExtraConfigs: keyscan.KeyScanExtraConfigs{
			DiodeDirection:   "col2row",
			PushThreshold:    1000,
			ReleaseThreshold: 400,
		},
		MatrixPins: keyscan.MatrixPins{
			Rows: []string{"D1", "D2"},
			Cols: []string{"D3", "D4"},
		},
	}
	
	// 初期化
	err := matrix.Init(config)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
	
	// スキャン
	isMatrixUpdate := matrix.Scan()
	
	// 結果の確認
	// 必要に応じてアサーションを追加
	
	// ファクトリー関数のテスト
	keyscan.RegisterKeyScanMatrix("new_type", func() keyscan.KeyScanMatrix { return keyscan.NewNewTypeKeyScan() })
	matrixFromFactory := keyscan.NewKeyScanMatrix("new_type")
	if matrixFromFactory == nil {
		t.Error("NewKeyScanMatrix returned nil for new_type")
	}
}
```

### 4. keyboard.jsonでの使用

新しいキースキャンタイプを使用するには、keyboard.jsonファイルの`keyscan`フィールドに新しいタイプの名前を指定します：

```json
{
    "name": "my_keyboard",
    "maintainer": "Your Name",
    "vendorId": "0xfeed",
    "productId": "0x6060",
    "keyscan": "new_type",
    "matrix": { "rows": 5, "cols": 6 },
    "keyscan_extra_configs": {
        "diode_direction": "col2row",
        "push_threshold": 1000,
        "release_threshold": 400
    },
    "hid": [
        "usb"
    ],
    "split": false,
    "matrix_pins": {
        "rows": ["D1", "D2", "D3", "D4", "D5"],
        "cols": ["D6", "D7", "D8", "D9", "D10", "D11"]
    }
}
```

## 注意事項

- 新しいキースキャンタイプは、`keyscan.KeyScanMatrix`インターフェースを実装する必要があります。
- ハードウェア固有の設定は、`KeyScanExtraConfigs`構造体に追加することができます。
- エラー処理を適切に行い、無効な設定や初期化エラーを明確に報告するようにしてください。
- テストを作成して、新しいキースキャンタイプが正しく動作することを確認してください。

## 既存のキースキャンタイプ

既存のキースキャンタイプの実装を参考にすることができます：

- `mx`: `keyscan/mx.go`
- `ec`: `keyscan/ec.go`

これらの実装は、新しいキースキャンタイプを作成する際の良いテンプレートとなります。