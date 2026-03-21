package device

import "ergo.services/ergo/gen"

type Device struct{}

func CreateDeviceApp() gen.ApplicationBehavior {
	return &Device{}
}

func (app *Device) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
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

func (app *Device) Start(_ gen.ApplicationMode) {}
func (app *Device) Terminate(_ error)           {}
