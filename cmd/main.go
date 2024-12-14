package main

import (
	"console-click-speed/pkg/game"
	"errors"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	text = "you must write it right"
)

func main() {
	gameGUI, err := game.NewGui(text)
	if err != nil {
		if errors.Is(err, gocui.ErrQuit) {
			gameGUI.Close()
			return
		}
		log.Panicln(err)
	}

	if err := gameGUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	defer gameGUI.Close()

	game.Start(gameGUI.Gui)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
