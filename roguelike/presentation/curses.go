package presentation

import (
	"github.com/rthornton128/goncurses"
)

const (
	WhiteFont  = 1
	RedFont    = 2
	GreenFont  = 3
	BlueFont   = 4
	YellowFont = 5
	CyanFont   = 6
)

func InitCurses() (*goncurses.Window, error) {
	stdscr, err := goncurses.Init()
	if err != nil {
		return nil, err
	}

	goncurses.Echo(false)
	goncurses.Cursor(0)
	stdscr.Keypad(true)
	goncurses.StartColor()

	if goncurses.HasColors() {
		goncurses.InitPair(WhiteFont, goncurses.C_WHITE, goncurses.C_BLACK)
		goncurses.InitPair(RedFont, goncurses.C_RED, goncurses.C_BLACK)
		goncurses.InitPair(GreenFont, goncurses.C_GREEN, goncurses.C_BLACK)
		goncurses.InitPair(BlueFont, goncurses.C_BLUE, goncurses.C_BLACK)
		goncurses.InitPair(YellowFont, goncurses.C_YELLOW, goncurses.C_BLACK)
		goncurses.InitPair(CyanFont, goncurses.C_CYAN, goncurses.C_BLACK)
	}

	return stdscr, nil
}

func GetInput(stdscr *goncurses.Window) int {
	key := stdscr.GetChar()
	charValue := int(key)
	return charValue & 0xFF
}
