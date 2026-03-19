package api

import "ergo.services/ergo/gen"

type Api struct{}

func CreateApiApp() gen.ApplicationBehavior {
	return &Api{}
}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *Api) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "api",
		Description: "API",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    "sup",
				Factory: newSup,
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *Api) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *Api) Terminate(reason error) {}
