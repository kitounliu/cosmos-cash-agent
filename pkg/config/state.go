package config

type Holding struct {
	Token  string
	Amount string
}

type State struct {
	Contacts      []string    `json:"contacts"`
	Balances      []Holding   `json:"balances"`
	Credentials   []string    `json:"credentials"`
	Notifications chan string `json:"-"`
}

func NewState() *State {
	return &State{
		Contacts:    make([]string, 0),
		Balances:    make([]Holding, 0),
		Credentials: make([]string, 0),
		Notifications: make(chan string),
	}
}
