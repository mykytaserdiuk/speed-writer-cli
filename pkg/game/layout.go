package game

import (
	"console-click-speed/pkg/gui"
	"console-click-speed/pkg/gui/editor"
	"console-click-speed/pkg/models"
	"fmt"

	"github.com/jroimartin/gocui"
)

func layout(text string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("game_1", 0, 0, maxX/2-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			v.Title = "WRITE TEXT FROM LEFT"
			editor := editor.NewManager(text, true, g)
			fmt.Fprint(v, editor.Text+"\n\n")
			v.SetCursor(0, 2)
			go func() {
				for {
					data := <-editor.KeyChan
					err := editor.UpdateTargetView(editor.G, data)
					if err != nil {
						fmt.Fprint(v, err.Error())
					}
				}
			}()
			v.Editor = editor

			if _, err = gui.SetCurrentViewOnTop(g, "game_1"); err != nil {
				return err
			}
		}

		if v, err := g.SetView("word_1", maxX/2+1, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Wrap = true
			v.Title = "TEXT"
			fmt.Fprint(v, models.GreenGBWhiteString+string(text[0])+models.NormalColorString+text[1:])
		}
		return nil
	}
}
