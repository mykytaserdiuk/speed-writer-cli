package gui

import "github.com/jroimartin/gocui"

var (
	viewArr = make([]string, 0)
	active  int
)

func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]
	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	// if nextIndex == 0 || nextIndex == 3 {
	// 	g.Cursor = true
	// } else {
	// 	g.Cursor = false
	// }

	active = nextIndex
	return nil
}
