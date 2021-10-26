package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type MenuWidget struct {
	name    string
	label   string
	buttons []*ButtonWidget
}

func NewMenuWidget(buttons ...int) *MenuWidget {

	var bw []*ButtonWidget
	bX := 0
	bY := 0
	for _, i := range buttons {
		b := NewButtonWidget(i, bX, bY, func(g *gocui.Gui, v *gocui.View) (err error) {
			state.Notifications <- fmt.Sprintf("clicked view %v, button %v: %v", v.Name(), i, label(i))
			_, err = g.SetCurrentView(v.Name())
			if err != nil {
				return
			}
			switch v.Name() {
			case fmt.Sprint(messagesMenuItem):
				NewMessagesView().Layout(g)
			case fmt.Sprint(balancesMenuItem):
				NewBalancesView().Layout(g)
			case fmt.Sprint(credentialsMenuItem):
				NewCredentialsView().Layout(g)
			default:
				NewDashboardView().Layout(g)
			}
			return
		})
		bw = append(bw, b)
		bX = b.x + b.w + 1
	}
	return &MenuWidget{name: name(header), buttons: bw}
}

func (w *MenuWidget) Layout(g *gocui.Gui) error {

	maxX, _ := g.Size()

	_, err := g.SetView(w.name, 0, 0, maxX-1, menuHeight)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	for _, b := range w.buttons {
		b.Layout(g)
	}
	return nil
}

type StatusbarWidget struct {
	name string
	val  string
}

func NewStatusbarWidget() *StatusbarWidget {
	return &StatusbarWidget{name: name(footer), val: "this is the status bar"}
}

func (w *StatusbarWidget) SetVal(val string) error {
	w.val = val
	return nil
}

func (w *StatusbarWidget) Val() string {
	return w.val
}

func (w *StatusbarWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	_, err := g.SetView(w.name, 0, maxY-footerHeight, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	// v.Clear()
	//rep := int(float64(maxX))
	//fmt.Fprint(v, strings.Repeat("▒", rep))\
	// fmt.Fprint(v, w.val)
	return nil
}
