package api

import "ergo.services/ergo/gen"

type Api struct{}

func CreateApiApp() gen.ApplicationBehavior { return &Api{} }

func (app *Api) Load(node gen.Node, args ...any) (spec gen.ApplicationSpec, err error) {
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

func (app *Api) Start(_ gen.ApplicationMode) {}
func (app *Api) Terminate(_ error)           {}
