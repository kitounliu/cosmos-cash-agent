package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ButtonWidget struct {
	name    string
	x, y    int
	w       int
	label   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewButtonWidget(id int, x, y int, handler func(g *gocui.Gui, v *gocui.View) error) *ButtonWidget {
	return &ButtonWidget{name: name(id), x: x, y: y, w: len(label(id)) + 1, label: label(id), handler: handler}
}

func (w *ButtonWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.handler); err != nil {
			return err
		}
		fmt.Fprint(v, w.label)
	}
	return nil
}

func NewListWidget(id int, x0, y0, x1, y1 int, items []string) *ListWidget {
	return &ListWidget{
		name:     name(id),
		title:    label(id),
		x0:       x0,
		y0:       y0,
		x1:       x1,
		y1:       y1,
		selected: 0,
		items:    items,
		handler:  nil,
	}
}

type ListWidget struct {
	name     string
	title    string
	x0, y0   int
	x1, y1   int
	selected int
	items    []string
	handler  func(g *gocui.Gui, v *gocui.View) error
}

func (w *ListWidget) Layout(g *gocui.Gui) error {

	v, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Title = w.title
	//if _, err := g.SetCurrentView(w.name); err != nil {
	//	return err
	//}
	if _, err := g.SetViewOnTop(w.name); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
		return err
	}
	v.Clear()
	for _, e := range w.items {
		fmt.Fprintln(v, e)
	}
	return nil
}

func NewTextInputWidget(id int, x0, y0, x1, y1 int) *TextInputWidget {
	return &TextInputWidget{
		name:    name(id),
		x0:      x0,
		y0:      y0,
		x1:      x1,
		y1:      y1,
		value:   "",
		handler: nil,
	}
}

type TextInputWidget struct {
	name    string
	x0, y0  int
	x1, y1  int
	value   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func (w *TextInputWidget) Layout(g *gocui.Gui) error {
	_, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.handler); err != nil {
			return err
		}
		g.Highlight = true
		g.SelFgColor = gocui.ColorRed
	}
	return nil
}
