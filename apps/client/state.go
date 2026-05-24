package client

import (
	"samovar/lib/change"
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

func (s *state) Init(_ ...any) error {
	s.data = st.New(st.Fields{
		st.FieldParams{"t_top", 'f', "t top", "°C"},
		st.FieldParams{"t_mid", 'f', "t middle", "°C"},
		st.FieldParams{"t_btm", 'f', "t bottom", "°C"},
	})

	s.Log().Debug("client.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("client.state received message: %#v", msg)

	switch v := msg.(type) {
	case change.Request:
		return s.updateState(v)
	default:
		err := fmt.Sprintf("client.state unexpected message: %#v", msg)
		return errors.New(err)
	}
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	s.Log().Debug("client.state got request: %#v", req)
	return gen.Atom("pong"), nil
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("client.state terminated: %s", reason)
}

func (s *state) updateState(req change.Request) error {
	s.Log().Debug("client.state change req: %#v", req)
	report, err := s.data.Change(req)

	if err == nil && report.Changed {
		report = report.AddFrom("client_state")

		if err = s.Send("client_producer", report); err == nil {
			s.Log().Debug(report.MakeLogString("client.state"))
		} else {
			return err
		}
	}

	return nil
}
