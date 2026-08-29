package session

import "ergo.services/ergo/gen"

type session struct{}

func CreateApp() gen.ApplicationBehavior { return &session{} }

func (app *session) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	cfg, err := loadConfig()

	spec.Name = "session_app"
	spec.Description = "SESSION"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "session_sup",
			Factory: newSup,
			Args:    []any{cfg},
		},
	}

	return spec, err
}

func (app *session) Start(_ gen.ApplicationMode) {}
func (app *session) Terminate(_ error)           {}
