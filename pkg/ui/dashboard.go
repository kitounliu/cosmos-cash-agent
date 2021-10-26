package ui

import (
	"github.com/jroimartin/gocui"
)

type DashboardView struct {
	name    string
	title   string
}

func NewDashboardView() *DashboardView {
	return &DashboardView{name: name(dashboardView),  title: label(dashboardView)}
}

func (w *DashboardView) Layout(g *gocui.Gui) error {
	var x0,y0,x1,y1 int
	x0,y0,x1,y1 = relativeXY(g, 0, 0, widthFull, heightFull)
	var widget gocui.Manager
	widget = NewListWidget(-1, x0, y0, x1, y1, []string{
		"balance 12,314.00€",
		"12 contacts",
		"41 credentials",
	})
	g.DeleteView("-1")
	return widget.Layout(g)
}