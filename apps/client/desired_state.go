package client

import (
	"samovar/lib/inter"
	st "samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type desiredState struct {
	act.Actor
	data *st.State
}

func newDesiredState() gen.ProcessBehavior {
	return &desiredState{}
}

func (ds *desiredState) Init(_ ...any) error {
	inter.RegisterActor(ds, "[(client.desired_state)]")

	ds.data = st.New(st.Fields{
		st.FieldParams{"collect_type", 's', "collect type", ""},
		st.FieldParams{"collect_value", 'i', "collect value", "%"},
	})

	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	if r, ok := req.(map[string]any); ok {
		if reports, err := s.data.BulkChangeFromMap(r, "client_desired_state"); err == nil {
			return inter.Send(s, reports, "reports", "client_producer")
		} else {
			return nil, err
		}
	} else {
		return s.data.HandleDataRequest(req)
	}
}
