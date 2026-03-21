package api

import "ergo.services/ergo/gen"

type api struct{}

func CreateApiApp() gen.ApplicationBehavior { return &api{} }

func (app *api) Load(node gen.Node, args ...any) (spec gen.ApplicationSpec, err error) {
	spec.Name = "api_app"
	spec.Description = "API"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "api_sup",
			Factory: newSup,
		},
	}

	return spec, err
}

func (app *api) Start(_ gen.ApplicationMode) {}
func (app *api) Terminate(_ error)           {}
