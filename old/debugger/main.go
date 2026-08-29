package main

import (
	"fmt"
	"time"

	"github.com/tgk-project/tgk"
	"github.com/tgk-project/tgk/di"
)

var tgkManager tgk.TGKManager

// like arduino main loop
func main() {
	setup()

	// デバッグモードで数回ループを実行
	fmt.Println("=== デバッグモード開始 ===")
	for i := 0; i < 10; i++ {
		fmt.Printf("--- ループ %d ---\n", i+1)
		loop()
		time.Sleep(100 * time.Millisecond) // 100ms待機
	}
	fmt.Println("=== デバッグモード終了 ===")
}

func setup() {
	fmt.Println("=== セットアップ開始 ===")
	tgkManager = di.InitializeTGKManager()

	// デバッグレベルを設定
	tgkManager.SetLogLevel(tgk.LogLevelDebug)

	tgkManager.SetLoadConfigName("keyboard.json")

	if err := tgkManager.Init(); err != nil {
		panic(err)
	}
	fmt.Println("=== セットアップ完了 ===")
}

func loop() {
	if err := tgkManager.Task(); err != nil {
		panic(err)
	}
}
