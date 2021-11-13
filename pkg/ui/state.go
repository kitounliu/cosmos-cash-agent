package ui

import "github.com/allinbits/cosmos-cash-agent/pkg/model"

type State struct {
	Contacts        map[string]model.Contact `json:"contacts"`
	SelectedContact int                      `json:"selected_contact"`
	Credentials     []string                 `json:"credentials"`
	Address         string                   `json:"address"`
}

func NewState() *State {
	return &State{
		Contacts:    make(map[string]model.Contact, 0),
		Credentials: make([]string, 0),
	}
}
