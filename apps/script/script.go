package script

import "ergo.services/ergo/gen"

type script struct{}

func CreateApp() gen.ApplicationBehavior { return &script{} }

func (app *script) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	spec.Name = "script_app"
	spec.Description = "SCRIPT"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "script_sup",
			Factory: newSup,
		},
	}

	return spec, err
}

func (app *script) Start(_ gen.ApplicationMode) {}
func (app *script) Terminate(_ error)           {}
