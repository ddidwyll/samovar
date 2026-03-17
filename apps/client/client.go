package client

import (
	"ergo.services/ergo/gen"
)

func CreateClient() gen.ApplicationBehavior {
	return &Client{}
}

type Client struct{}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *Client) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "client",
		Description: "description of this application",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    "clientsup",
				Factory: factory_ClientSup,
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *Client) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *Client) Terminate(reason error) {}
