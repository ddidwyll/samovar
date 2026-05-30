package device

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type sup struct{ act.Supervisor }

func newSup() gen.ProcessBehavior { return &sup{} }

func (s *sup) Init(_ ...any) (spec act.SupervisorSpec, err error) {
	spec.Type = act.SupervisorTypeOneForOne
	spec.Restart.Strategy = act.SupervisorStrategyTransient
	spec.Restart.Intensity = 2
	spec.Restart.Period = 5

	spec.Children = []act.SupervisorChildSpec{
		{
			Name:    "device_change_log",
			Factory: newChangeLog,
		},
		{
			Name:    "device_raw_state",
			Factory: newRawState,
		},
		{
			Name:    "device_state",
			Factory: newState,
		},
		{
		  Name: "device_calc",
		  Factory: newCalc,
		},
	}

	return spec, err
}
