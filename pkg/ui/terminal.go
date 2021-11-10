package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type TerminalWidget struct {
	name    string
	val     string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewTerminalWidget() *TerminalWidget {
	return &TerminalWidget{
		name:    name(footer),
		val:     "this is the status bar",
		handler: TerminalCommandFunc,
	}
}

func (w *TerminalWidget) SetVal(val string) error {
	w.val = val
	return nil
}

func (w *TerminalWidget) Val() string {
	return w.val
}

func (w *TerminalWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	_, err := g.SetView(w.name, 0, maxY-footerHeight, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
		fmt.Println("asdasdasd")
		return err
	}
	//fmt.Println(w.name, w.val)
	//w.Clear()
	//rep := int(float64(maxX))
	//fmt.Fprint(v, strings.Repeat("▒", rep))\
	// fmt.Fprint(v, w.val)
	return nil
}
