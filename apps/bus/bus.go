package bus

import "ergo.services/ergo/gen"

type bus struct{}

func CreateApp() gen.ApplicationBehavior {
	return &bus{}
}

func (app *bus) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	spec.Name = "bus_app"
	spec.Description = "BUS"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "bus_sup",
			Factory: newSup,
		},
	}

	return spec, err
}

func (app *bus) Start(_ gen.ApplicationMode) {}
func (app *bus) Terminate(_ error)           {}
