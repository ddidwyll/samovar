package client

import (
  "samovar/lib/state"
  st "samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type state struct{
  act.Actor
  data *st.State
}

func newState() gen.ProcessBehavior {
  return &state{}
}

func (cs *state) Init(_ ...any) error {
	cs.data = st.New(st.Fields{
		st.FieldParams{"t_top", 'f', "t top", "°C"},
		st.FieldParams{"t_mid", 'f', "t middle", "°C"},
		st.FieldParams{"t_btm", 'f', "t bottom", "°C"},
	})

	cs.Log().Debug("client.state started (%s)", cs.Name())
	return nil
}

func (cs *state) HandleMessage(_ gen.PID, msg any) error {
	cs.Log().Debug("client.state received message: %#v", msg)
	return nil
}

func (cs *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	cs.Log().Debug("client.state got request: %#v", req)
	return gen.Atom("pong"), nil
}

func (cs *state) Terminate(reason error) {
	cs.Log().Debug("terminated: %s", reason)
}
