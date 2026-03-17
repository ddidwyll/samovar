package device

import (
	"ergo.services/ergo/gen"
)

func CreateDevice() gen.ApplicationBehavior {
	return &Device{}
}

type Device struct{}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *Device) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "device",
		Description: "description of this application",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    "devicesup",
				Factory: factory_DeviceSup,
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *Device) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *Device) Terminate(reason error) {}
