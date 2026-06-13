package client

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	st "samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type state struct {
	act.Actor
	data *st.State
}

func newState() gen.ProcessBehavior {
	return &state{}
}

var stateFields = st.Fields{
	st.FieldParams{"t_top", 'f', "t top", "°C"},
	st.FieldParams{"t_mid", 'f', "t middle", "°C"},
	st.FieldParams{"t_btm", 'f', "t bottom", "°C"},
	st.FieldParams{"power", 'i', "power", "W"},
	st.FieldParams{"power_diff", 'i', "power diff", "%"},
	st.FieldParams{"collect", 's', "collect", "%"},
	st.FieldParams{"press", 'f', "press", "mm"},
	st.FieldParams{"devices", 'a', "devices", ""},
	st.FieldParams{"device_id", 's', "device id", ""},
	st.FieldParams{"error", 's', "error", ""},
	st.FieldParams{"started", 's', "started", ""},
}

func (s *state) Init(_ ...any) error {
	inter.RegisterActor(s, "[(client.state)]")
	s.data = st.New(stateFields)

	s.Log().Debug("client.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("client.state received message: %#v", msg)

	switch r := msg.(type) {
	case []change.Request:
		return s.updateState(r...)
	case change.Request:
		return s.updateState(r)
	default:
		err := fmt.Sprintf("client.state unexpected message: %#v", msg)
		return errors.New(err)
	}
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.data.HandleDataRequest(req)
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("client.state terminated: %s", reason)
}

func (s *state) updateState(reqs ...change.Request) error {
	s.Log().Debug("client.state change requests: %#v", reqs)

	if reports, err := s.data.BulkChange(reqs, "client_state"); err != nil {
		return err
	} else {
		return inter.Send(s, reports, "reports", "client_producer")
	}
}

func stateEntries(p gen.Process) (st.Entries, error) {
	if kv, err := inter.Call(p, nil, "state", "client_state"); err != nil {
		return nil, err
	} else {
		if result, ok := kv.(st.Entries); !ok {
			return nil, errors.New("Unexpected error")
		} else {
			return result, nil
		}
	}
}
