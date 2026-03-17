package script

import (
	"ergo.services/ergo/gen"
)

func CreateScript() gen.ApplicationBehavior {
	return &Script{}
}

type Script struct{}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *Script) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "script",
		Description: "description of this application",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    "scriptsup",
				Factory: factory_ScriptSup,
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *Script) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *Script) Terminate(reason error) {}
