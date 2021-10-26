package ui

import (
	"github.com/jroimartin/gocui"
)

type MessagesView struct {
	name  string
	title string
}

func NewMessagesView() *MessagesView {
	return &MessagesView{name: name(messagesView), title: label(messagesView)}
}

func (w *MessagesView) Layout(g *gocui.Gui) error {
	var x0, y0, x1, y1 int
	var widget gocui.Manager
	x0, y0, x1, y1 = relativeXY(g, 0, 0, widthQuarter, heightFull)
	cl := NewListWidget(contactsPanel, x0, y0, x1, y1, state.Contacts)
	cl.Layout(g)

	x0, y0, x1, y1 = relativeXY(g, x1, 0, widthFull, 0.9)
	widget = NewListWidget(chatPanel, x0, y0, x1, y1, []string{})
	if err := widget.Layout(g); err != nil {
		return err
	}

	x0, y0, x1, y1 = relativeXY(g, x0, y1, widthFull, heightFull)
	widget = NewTextInputWidget(chatSendPanel, x0, y0, x1, y1)
	if err := widget.Layout(g); err != nil {
		return err
	}

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
