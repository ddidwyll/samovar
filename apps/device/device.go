package device

import "ergo.services/ergo/gen"

type device struct{}

func CreateDeviceApp() gen.ApplicationBehavior {
	return &device{}
}

func (app *device) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	spec.Name = "device_app"
	spec.Description = "DEVICE"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "device_sup",
			Factory: newSup,
		},
	}

	return spec, err
}

func (app *device) Start(_ gen.ApplicationMode) {}
func (app *device) Terminate(_ error)           {}
