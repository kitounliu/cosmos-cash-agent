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
	"github.com/allinbits/cosmos-cash-agent/pkg/helpers"
)

var (
	appCfg *config.EdgeConfigSchema
	// ui data binding
	userCommand = binding.NewString()
	// messages tab command
	messages = binding.NewStringList()
	// balances
	balances             = binding.NewStringList()
	balancesChainOfTrust = binding.NewString()
	// credentials
	// TODO need to have separate stuff for public and private credentials
	credentials    = binding.NewStringList()
	credentialData = binding.NewString()
	// contacts
	contacts    = binding.NewStringList()
	contactData = binding.NewString()
	// logs
	logData = binding.NewString()
)

func Render(cfg *config.EdgeConfigSchema) {

	appCfg = cfg

	myApp := app.New()
	myWindow := myApp.NewWindow(cfg.ControllerName)

	// main content
	tabs := container.NewAppTabs(
		getMessagesTab(),
		getCredentialsTab(),
		getBalancesTab(),
		getDashboardTab(),
		getLogTab(),
	)

	myWindow.SetContent(
		container.NewMax(
			tabs,
			//footer,
		),
	)

	myWindow.SetOnClosed(func() {
		// write the state on file
		appState, _ := config.GetAppData("state.json")
		helpers.WriteJson(appState, cfg.RuntimeState)
	})

	// run the dispatcher that updates the ui
	go dispatcher(cfg.RuntimeMsgs.Notification)
	// lanuch the app
	myWindow.Resize(fyne.NewSize(940, 660))
	myWindow.ShowAndRun()
}

func getMessagesTab() *container.TabItem {

	list := widget.NewListWithData(
		contacts,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(di.(binding.String))
		},
	)
	list.OnSelected = contactSelected

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

	body := container.NewHSplit(list, rightPanel)
	main := container.New(layout.NewMaxLayout(), body)

	return container.NewTabItem("Messages", main)
}

func getDashboardTab() *container.TabItem {

	main := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Wallet owner: \n%s", appCfg.ControllerName)),
		widget.NewLabel(fmt.Sprintf("Wallet DID ID: \ndid:cosmos:net:%s:%s", appCfg.ChainID, appCfg.ControllerDidID)),
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
			o.(*widget.Label).Bind(di.(binding.String))
		},
	)
	list.OnSelected = credentialsSelected

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
			o.(*widget.Label).Bind(di.(binding.String))
		},
	)

	list.OnUnselected = balancesSelected

	msgPanel := widget.NewLabelWithData(balancesChainOfTrust)
	msgScroll := container.NewScroll(msgPanel)

	body := container.NewHSplit(list, msgScroll)
	main := container.New(layout.NewMaxLayout(), body)

	return container.NewTabItem("Balances", main)
}

func getLogTab() *container.TabItem {
	msgPanel := widget.NewLabelWithData(logData)
	main := container.NewScroll(msgPanel)
	return container.NewTabItem("Logs", main)
}
