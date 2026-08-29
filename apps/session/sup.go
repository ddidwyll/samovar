package session

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type sup struct{ act.Supervisor }

func newSup() gen.ProcessBehavior { return &sup{} }

func (s *sup) Init(args ...any) (spec act.SupervisorSpec, err error) {
	spec.Type = act.SupervisorTypeOneForOne
	spec.Restart.Strategy = act.SupervisorStrategyTransient
	spec.Restart.Intensity = 2
	spec.Restart.Period = 5

	spec.Children = []act.SupervisorChildSpec{
		{
			Name:    "session_state",
			Factory: newState,
		},
		{
			Name:    "session_desired_state",
			Factory: newDesiredState,
		},
		{
			Name:    "session_repo",
			Factory: newRepo,
		},
		{
			Name:    "session_calc",
			Factory: newCalc,
			Args:    args,
		},
	}

	return spec, err
}
