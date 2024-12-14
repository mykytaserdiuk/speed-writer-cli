package models

const (
	NormalColorString  = "\033[0m"
	RedColorString     = "\033[31m"
	RedGBWhiteString   = "\033[41;37m"
	GreenGBWhiteString = "\033[42;37m"
)

func StyledRune(style string, text rune) string {
	st := style + string(text) + NormalColorString
	return st
}
