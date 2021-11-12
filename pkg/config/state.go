package config

type Holding struct {
	Token  string
	Amount string
}

type State struct {
	Contacts    map[string]Contact `json:"contacts"`
	Credentials []string           `json:"credentials"`
	Address     string             `json:"address"`
}

func NewState() *State {
	return &State{
		Contacts:    make(map[string]Contact, 0),
		Credentials: make([]string, 0),
	}
}

type Contact struct {
	DID     string
	Address string
	Name    string
}
