package editor

import (
	"console-click-speed/pkg/models"
	"console-click-speed/pkg/save"
	"console-click-speed/pkg/utils"
	"fmt"
	"strings"
	"unicode"

	"github.com/jroimartin/gocui"
)

func NewManager(text string, spaceMode bool, g *gocui.Gui) *Manager {
	return &Manager{
		Text:       text,
		spaceMode:  spaceMode,
		KeyChan:    make(chan Data),
		G:          g,
		errorCount: 0,
	}
}

type Manager struct {
	Text       string
	currentPos int
	// if false = skip a 'space '
	spaceMode  bool
	KeyChan    chan Data
	errorCount int
	G          *gocui.Gui
	resultText string
}

type Data struct {
	v   *gocui.View
	key gocui.Key
	ch  rune
	mod gocui.Modifier
}

func (ve *Manager) scoreErrorPercent() float32 {
	le := float32(len(ve.Text))
	if le == 0 {
		return 0
	}
	correctCount := le - float32(ve.errorCount)
	return (correctCount / le) * 100
}

func (ve *Manager) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if !utils.CheckAllowingKey(key) {
		return
	}

	ve.KeyChan <- Data{v, key, ch, mod}
	shouldSkipValidation := false
	newString := ""
	v.MoveCursor(+1, 0, true)
	if ve.currentPos == len(ve.Text) {
		v.Clear()
		persent := ve.scoreErrorPercent()
		fmt.Fprint(v, persent)
		err := save.SaveInResult(ve.resultText, 1, persent)
		if err != nil {
			panic(err)
		}
		return
	}

	whatMustBe := rune(ve.Text[ve.currentPos])

	if ve.spaceMode {
		if unicode.IsSpace(whatMustBe) {
			if key == gocui.KeySpace {
				shouldSkipValidation = true
				newString = " "
			} else {
				newString = models.StyledRune(models.RedGBWhiteString, ch)
				ve.errorCount++
				shouldSkipValidation = true
			}
		} else if key == gocui.KeySpace {
			newString = models.StyledRune(models.RedGBWhiteString, '_')
			ve.errorCount++
			shouldSkipValidation = true
		}
	}

	if !ve.spaceMode && whatMustBe == ' ' {
		ve.currentPos++
		whatMustBe = rune(ve.Text[ve.currentPos])
		v.MoveCursor(+1, 0, true)
	}

	if !shouldSkipValidation {
		if whatMustBe == ch {
			newString = models.StyledRune(models.GreenGBWhiteString, ch)
		} else {
			ve.errorCount++
			newString += models.StyledRune(models.RedColorString, ch)
		}
	}

	fmt.Fprint(v, newString)
	ve.resultText += newString
	ve.currentPos++
}

func (ve *Manager) UpdateTargetView(g *gocui.Gui, data Data) error {
	v, err := g.View("word_1")
	if err != nil {
		return err
	}
	v.Clear()
	x, y := data.v.Cursor()
	for i, st := range strings.Split(ve.Text, "\n") {
		for j, ch := range st {
			if y-2 == i && x == j {
				v.Write([]byte(models.StyledRune(models.GreenGBWhiteString, ch)))
			} else {
				v.Write([]byte(string(ch)))
			}
		}
		v.Write([]byte("\n\n"))
	}
	return nil
}
