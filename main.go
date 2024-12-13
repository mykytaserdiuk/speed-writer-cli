package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	normalColorString  = "\033[0m"
	redColorString     = "\033[31m"
	redGBWhiteString   = "\033[41;37m"
	greenGBWhiteString = "\033[42;37m"
)

var (
	viewArr = []string{"game_1", "word_1"}
	active  = 0
	text    = "you must write it right"
)

type Editor struct {
	currentPos int
	// if false = skip a 'space '
	spaceMode  bool
	keyChan    chan Data
	errorCount int
	g          *gocui.Gui
	resultText string
}

type Data struct {
	v   *gocui.View
	key gocui.Key
	ch  rune
	mod gocui.Modifier
}

func (ve *Editor) scoreErrorPersent() float32 {
	le := float32(len(text))
	return float32(len(text)-ve.errorCount) / le * 100
}

func (ve *Editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if !CheckAllowingKey(key) {
		return
	}

	ve.keyChan <- Data{v, key, ch, mod}
	skipCheck := false
	newString := ""
	v.MoveCursor(+1, 0, true)
	if ve.currentPos == len(text) {
		v.Clear()
		persent := ve.scoreErrorPersent()
		fmt.Fprint(v, persent)
		err := SaveInResult(ve.resultText, 1, persent)
		if err != nil {
			panic(err)
		}
		return
	}
	whatMustBe := rune(text[ve.currentPos])

	if ve.spaceMode && whatMustBe == ' ' && key == gocui.KeySpace {
		newString = " "
	} else if ve.spaceMode && whatMustBe == ' ' && key != gocui.KeySpace {
		newString = redGBWhiteString + string(ch) + normalColorString
		ve.errorCount++
		skipCheck = true
	} else if !ve.spaceMode && whatMustBe == ' ' {
		ve.currentPos++
		whatMustBe = rune(text[ve.currentPos])
		v.MoveCursor(+1, 0, true)
	} else if ve.spaceMode && key == gocui.KeySpace && whatMustBe != ' ' {
		newString = redGBWhiteString + string("_") + normalColorString
		ve.errorCount++
		skipCheck = true
	}
	if !skipCheck {
		if whatMustBe == ch {
			newString = greenGBWhiteString + string(ch) + normalColorString
		} else {
			ve.errorCount++
			newString += redColorString + string(ch) + normalColorString
		}
	}

	fmt.Fprint(v, newString)
	ve.resultText += newString
	ve.currentPos++
}

func updateTargetView(g *gocui.Gui, data Data) error {
	v, err := g.View("word_1")
	if err != nil {
		return err
	}
	v.Clear()
	x, y := data.v.Cursor()
	for i, st := range strings.Split(text, "\n") {
		for j, ch := range st {
			if y-2 == i && x == j {
				v.Write([]byte(greenGBWhiteString + string(ch) + normalColorString))
			} else {
				v.Write([]byte(string(ch)))
			}
		}
		v.Write([]byte("\n\n"))
	}
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
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

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("game_1", 0, 0, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "WRITE TEXT FROM LEFT"
		fmt.Fprint(v, text+"\n\n")
		v.SetCursor(0, 2)
		editor := &Editor{spaceMode: true, keyChan: make(chan Data), g: g}
		go func() {
			for {
				data := <-editor.keyChan
				err := updateTargetView(editor.g, data)
				if err != nil {
					fmt.Fprint(v, err.Error())
				}
			}
		}()
		v.Editor = editor

		if _, err = setCurrentViewOnTop(g, "game_1"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("word_1", maxX/2+1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = "TEXT"
		fmt.Fprint(v, greenGBWhiteString+string(text[0])+normalColorString+text[1:])
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
