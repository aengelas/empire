package relay

import (
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

var (
	DefaultSessionGenerator = func() string { return uuid.New() }

	// DefaultOptions is a default Options instance that can be passed when
	// intializing a new Relay.
	DefaultOptions = Options{SessionGenerator: DefaultSessionGenerator}
)

type Options struct {
	Host             string
	SessionGenerator func() string
}

// Container represents a docker container to run.
type Container struct {
	Image     string            `json:"image"`
	Name      string            `json:"name"`
	Command   string            `json:"command"`
	Status    string            `json:"status"`
	Env       map[string]string `json:"env"`
	Attach    bool              `json:"attach"`
	AttachURL string            `json:"attach_url"`
}

type Relay struct {
	sync.Mutex

	// The rendezvous host
	Host string

	genSessionId func() string
	sessions     map[string]bool
}

// New returns a new Relay instance.
func New(options Options) *Relay {
	sg := options.SessionGenerator
	if sg == nil {
		sg = DefaultSessionGenerator
	}

	return &Relay{
		Host:         options.Host,
		genSessionId: sg,
		sessions:     map[string]bool{},
	}
}

func (r *Relay) NewSession() string {
	r.Lock()
	defer r.Unlock()
	id := r.genSessionId()
	r.sessions[id] = true
	return id
}

func (r *Relay) CreateContainer(c *Container) error {
	// docker pull
	// docker run
	return nil
}
