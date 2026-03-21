package device

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Sup struct{ act.Supervisor }

func newSup() gen.ProcessBehavior { return &Sup{} }

func (sup *Sup) Init(_ ...any) (spec act.SupervisorSpec, err error) {
	spec.Type = act.SupervisorTypeOneForOne
	spec.Restart.Strategy = act.SupervisorStrategyTransient
	spec.Restart.Intensity = 2
	spec.Restart.Period = 5

	spec.Children = []act.SupervisorChildSpec{
		{
			Name:    "device_raw_state",
			Factory: newRawState,
		},
	}

	return spec, err
}
