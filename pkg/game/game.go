package game

import (
	"console-click-speed/pkg/gui"
	"log"

	"github.com/jroimartin/gocui"
)

func NewGui(text string) (*gui.GUI, error) {
	gameGUI, err := gui.NewGui(gocui.OutputNormal, text)
	if err != nil {
		return nil, err
	}

	if err := gameGUI.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.NextView); err != nil {
		log.Panicln(err)
	}

	gameGUI.Highlight = true
	gameGUI.Cursor = true
	gameGUI.SelFgColor = gocui.ColorGreen

	gameGUI.SetManagerFunc(layout(text))

	return gameGUI, nil
}

func Start(gameGUI *gocui.Gui) error {
	if err := gameGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
