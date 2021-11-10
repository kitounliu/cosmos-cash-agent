package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type CommandFunc func(*gocui.Gui, *gocui.View) error

func TerminalCommandFunc(*gocui.Gui, *gocui.View) error {
	fmt.Println("hope")
	return nil
}
