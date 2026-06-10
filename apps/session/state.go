package session

import (
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

var stateFields = st.Fields{
	st.FieldParams{"device_id", 's', "device id", ""},
}

func (s *state) Init(_ ...any) error {
	inter.RegisterActor(s, "[(session.state)]")
	s.data = st.New(stateFields)

	s.Log().Debug("session.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.data.HandleReq(req)
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("session.state terminated: %s", reason)
}
