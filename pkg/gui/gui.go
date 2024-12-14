package gui

import (
	"console-click-speed/pkg/gui/editor"
	"github.com/jroimartin/gocui"
)

type GUI struct {
	Manager *editor.Manager
	*gocui.Gui
}

func NewGui(mode gocui.OutputMode, text string) (*GUI, error) {
	g, err := gocui.NewGui(mode)
	if err != nil {
		return nil, err
	}
	manager := editor.NewManager(text, true, g)
	return &GUI{
		manager,
		g,
	}, nil
}
