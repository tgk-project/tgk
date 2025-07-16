package main

import (
	"github.com/Diwamoto/tgk"
	"github.com/Diwamoto/tgk/di"
)

var tgkManager tgk.TGKManager

// like arduino main loop
func main() {
	setup()
	// for {
	// 	loop()
	// }
}

func setup() {
	tgkManager = di.InitializeTGKManager()
	tgkManager.SetLoadConfigName("keyboard.json")

	if err := tgkManager.Init(); err != nil {
		panic(err)
	}

}

// func loop() {
// 	if err := tgkManager.Task(); err != nil {
// 		panic(err)
// 	}
// }
