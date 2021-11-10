package ui

import (
	"fmt"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/jroimartin/gocui"
	"log"
)

const (
	menuHeight   = 3
	footerHeight = 3
	leftMargin   = 0
	rightMargin  = 0
)

const (
	heightFull    = 1.0
	widthFull     = 1.0
	heightHalf    = 0.5
	widthHalf     = 0.5
	heightQuarter = 0.25
	widthQuarter  = 0.25
)

const (
	// main views
	mainView = iota
	credentialsView
	messagesView
	dashboardView
	balancesView
	// header/footer
	header
	footer
	// menu buttons
	credentialsMenuItem
	messagesMenuItem
	dashboardMenuItem
	balancesMenuItem
	// messages panels
	contactsPanel
	chatPanel
	chatSendPanel
	// credentials panels
	credentialListPanel
	credentialDetailPanel
	// balances panels
	txHistoryPanel
	tokensPanel
)

var (
	labels = map[int]string{
		credentialsView: "Credentials",
		messagesView:    "Messages",
		dashboardView:   "Dashboard",
		balancesView:    "Balances",
		// menu
		credentialsMenuItem: "Credentials",
		messagesMenuItem:    "Messages",
		dashboardMenuItem:   "Dashboard",
		balancesMenuItem:    "Balances",
		// messages panels
		contactsPanel: "Contacts",
		chatPanel:     "Chat",
		chatSendPanel: "Message",
		// credentials panels
		credentialListPanel:   "Credential list",
		credentialDetailPanel: "Credential detail",
		txHistoryPanel:        "Transactions History",
		tokensPanel:           "Stable Coins",
	}
)

// l return the id and the label for a component
func l(id int) (string, string) {
	return name(id), label(id)
}

func name(id int) string {
	return fmt.Sprint(id)
}

func label(id int) string {
	return labels[id]
}

var (
	state *config.State
	hub   *config.MsgHub
)

func Render(appState *config.State, msgHub *config.MsgHub) {
	state = appState
	hub = msgHub
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Mouse = true

	menu := NewMenuWidget(
		dashboardMenuItem,
		messagesMenuItem,
		balancesMenuItem,
		credentialsMenuItem,
	)
	statuBar := NewTerminalWidget()

	// main views
	dash := NewDashboardView()

	g.SetManager(menu, statuBar, dash)

	go func() {
		for {
			n := <-hub.Notification
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(name(footer))
				if err != nil {
					return err
				}
				v.Clear()
				fmt.Fprintln(v, n)
				return nil
			})
		}
	}()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err.Error())
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// relativeXY calculate the correct x/y in a container
func relativeXY(g *gocui.Gui, x0, y0 int, w, h float32) (_x0 int, _y0 int, _x1 int, _y1 int) {
	mx, my := g.Size()

	_x0 = x0 + leftMargin
	_x1 = _x0 + int(float32(mx-(_x0+rightMargin))*w)
	_x1--

	_y0 = y0 + menuHeight
	_y1 = _y0 + int(float32(my-(_y0+footerHeight))*h)
	_y1--

	return
}
