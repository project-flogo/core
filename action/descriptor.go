package action

// Descriptor is the descriptor for the Action
type Descriptor struct {
	ID      string `json:"ref"`
	Version string `json:"version"`
}

func NewInfo(async bool, passthru bool) *Info {
	return &Info{async: async, passthru: passthru}
}

type Info struct {
	async    bool
	passthru bool
	id       string
}

func (i *Info) Id() string {
	return i.id
}

func (i *Info) Async() bool {
	return i.async
}

func (i *Info) Passthru() bool {
	return i.passthru
}
