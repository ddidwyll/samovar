package device

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

	s.Log().Debug("device.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("device.state received message: %#v", msg)

	switch req := msg.(type) {
	case change.Request:
		switch req.LastFrom() {
		case "device_raw_state":
			return s.applyFromRaw(req)
		default:
			return nil
		}
	default:
		err := fmt.Sprintf("device.state unexpected message: %#v", msg)
		return errors.New(err)
	}
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	s.Log().Debug("device.state got request: %#v", req)
	return gen.Atom("pong"), nil
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("device.state terminated: %s", reason)
}

func (s *state) updateState(req change.Request) error {
	s.Log().Debug("device.state change req: %#v", req)
	report, err := s.data.Change(req)

	if err == nil && report.Changed {
		report = report.AddFrom("device_state")
		// err = s.Send("device_producer", report)
		s.Log().Info(report.MakeLogString("device.state"))
	}

	return err
}

func (s *state) applyFromRaw(req change.Request) error {
	if req, mapped := mapFromRaw(req); mapped {
		return s.updateState(req)
	}
	return nil
}
