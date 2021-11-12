package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
)

var (
	appCfg config.EdgeConfigSchema
	// ui data binding
	userCommand = binding.NewString()
	// messages tab command
	messages = binding.NewStringList()
	// balances
	balances             = binding.NewStringList()
	balancesChainOfTrust = binding.NewString()
	// credentials
	credentials          = binding.NewStringList()
	credentialData       = binding.NewString()
)

func Render(cfg config.EdgeConfigSchema) {

	appCfg = cfg

	myApp := app.New()
	myWindow := myApp.NewWindow(cfg.ControllerName)

	tabs := container.NewAppTabs(
		getMessagesTab(),
		getCredentialsTab(),
		getBalancesTab(),
		getDashboardTab(),
	)

	myWindow.SetContent(
		container.NewMax(
			tabs,
			//footer,
		),
	)

	myWindow.SetOnClosed(func() {
		// write the state on file
	})

	myWindow.Resize(fyne.NewSize(940, 660))
	myWindow.ShowAndRun()
}

func getMessagesTab() *container.TabItem {
	contacts := []string{
		"alice",
		"bob",
		"emti",
		"whatever",
	}

	contactList := widget.NewList(
		func() int {
			return len(contacts)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(contacts[id])
		})

	msgPanel := container.NewVBox()
	msgScroll := container.NewScroll(msgPanel)

	// footer stuff
	rightPanel := container.NewBorder(
		nil,
		container.NewBorder(
			nil,
			nil,
			nil,
			widget.NewButtonWithIcon("", theme.MailSendIcon(), executeCmd),
			widget.NewEntryWithData(userCommand),
		),
		nil,
		nil,
		msgScroll,
	)

	body := container.NewHSplit(contactList, rightPanel)
	main := container.New(layout.NewMaxLayout(), body)

	return container.NewTabItem("Messages", main)
}

func getDashboardTab() *container.TabItem {

	main := container.NewVBox(

		widget.NewLabel(fmt.Sprintf("Wallet owner: %s", appCfg.ControllerName)),
		widget.NewLabel(fmt.Sprintf("Wallet DID ID: did:cosmos:net:%s:%s", appCfg.ChainID, appCfg.ControllerDidID)),
	)

	return container.NewTabItem("Dashboard", main)
}

func getCredentialsTab() *container.TabItem {
	list := widget.NewListWithData(
		credentials,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {

			// o.(*widget.Label).SetText(id.)
		},
	)


	msgPanel := widget.NewLabelWithData(credentialData)
	msgScroll := container.NewScroll(msgPanel)

	body := container.NewHSplit(list, msgScroll)
	main := container.New(layout.NewMaxLayout(), body)

	return container.NewTabItem("Credentials", main)
}

func getBalancesTab() *container.TabItem {

	list := widget.NewListWithData(
		balances,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {

			// o.(*widget.Label).SetText(id.)
		},
	)

	list.OnUnselected = balancesClick

	msgPanel := widget.NewLabelWithData(balancesChainOfTrust)
	msgScroll := container.NewScroll(msgPanel)

	body := container.NewHSplit(list, msgScroll)
	main := container.New(layout.NewMaxLayout(), body)

	return container.NewTabItem("Balances", main)
}
