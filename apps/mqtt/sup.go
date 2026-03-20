package mqtt

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Sup struct {
	act.Supervisor
}

func newSup() gen.ProcessBehavior {
	return &Sup{}
}

func (sup *Sup) Init(args ...any) (act.SupervisorSpec, error) {
	var spec act.SupervisorSpec

	spec.Type = act.SupervisorTypeOneForOne
	spec.Restart.Strategy = act.SupervisorStrategyTransient
	spec.Restart.Intensity = 2
	spec.Restart.Period = 5

	spec.Children = []act.SupervisorChildSpec{
		{
			Name:    "server",
			Factory: newServer,
			Args:    args,
		},
	}

	return spec, nil
}
