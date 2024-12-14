package utils

import "github.com/jroimartin/gocui"

var (
	NotAllowedKeys = []gocui.Key{gocui.KeyArrowDown, gocui.KeyArrowUp, gocui.KeyArrowLeft, gocui.KeyArrowRight, gocui.KeyDelete}
)

func CheckAllowingKey(key gocui.Key) bool {
	for _, k := range NotAllowedKeys {
		if k == key {
			return false
		}
	}
	return true
}
