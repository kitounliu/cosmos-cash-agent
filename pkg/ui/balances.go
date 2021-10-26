package ui

import (
	"github.com/jroimartin/gocui"
)

type BalancesView struct {
	name  string
	title string
}

func NewBalancesView() *BalancesView {
	return &BalancesView{name: name(balancesView), title: label(balancesView)}
}

func (w *BalancesView) Layout(g *gocui.Gui) error {
	var x0, y0, x1, y1 int
	x0, y0, x1, y1 = relativeXY(g, 0, 0, widthHalf, heightFull)
	cl := NewListWidget(tokensPanel, x0, y0, x1, y1, []string{
		"sEURO: 301.00€",
	})
	cl.Layout(g)

	x0, y0, x1, y1 = relativeXY(g, x1, 0, widthFull, heightFull)
	ml := NewListWidget(txHistoryPanel, x0, y0, x1, y1, []string{})
	ml.Layout(g)

	//v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	//if err != nil {
	//	if err != gocui.ErrUnknownView {
	//		return err
	//	}
	//	if _, err := g.SetCurrentView(w.name); err != nil {
	//		return err
	//	}
	//	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
	//		return err
	//	}
	//	fmt.Fprint(v, w.label)
	//}
	return nil
}
