package main

import (
	_ "embed"

	"github.com/tgk-project/tgk"
)

//go:embed keyboard.json
var keyboardJsonData []byte

func main() {
	keyboard := tgk.NewKeyboard(keyboardJsonData)
	keyboard.Start()
}
