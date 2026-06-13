package device

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	st "samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type state struct {
	act.Actor
	data *st.State
}

func newState() gen.ProcessBehavior {
	return &state{}
}

func (s *state) Init(_ ...any) error {
	inter.RegisterActor(s, "[(device.state)]")

	s.data = st.New(st.Fields{
		st.FieldParams{"t_top", 'f', "t top", "°C"},
		st.FieldParams{"t_mid", 'f', "t middle", "°C"},
		st.FieldParams{"t_btm", 'f', "t bottom", "°C"},
		st.FieldParams{"power", 'i', "power", "W"},
		st.FieldParams{"power_diff", 'i', "power diff", "%"},
		st.FieldParams{"collect", 's', "collect", "%"},
		st.FieldParams{"collect_value", 'f', "collect value", "%"},
		st.FieldParams{"collect_mode", 's', "collect mode", ""},
		st.FieldParams{"press", 'f', "press", "mm"},
	})

	s.Log().Debug("device.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleMessage(_ gen.PID, req any) error {
  reporter := func(reports ...change.Report) error {
    return inter.Send(s, reports, "reports", "device_producer")
  }
  return s.data.HandleChangeRequests(req, "device_state", reporter)
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.data.HandleDataRequest(req)
}
