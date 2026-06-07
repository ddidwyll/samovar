package client

import "ergo.services/ergo/gen"

type client struct{}

func CreateApp() gen.ApplicationBehavior { return &client{} }

func (app *client) Load(node gen.Node, args ...any) (spec gen.ApplicationSpec, err error) {
	cfg, err := loadConfig()

	spec.Name = "client_app"
	spec.Description = "CLIENT"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "client_sup",
			Factory: newSup,
			Args:    []any{cfg},
		},
	}

	return spec, err
}

func (app *client) Start(_ gen.ApplicationMode) {}
func (app *client) Terminate(_ error)           {}
