package device

import (
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
	ds.data = st.New(st.Fields{
		st.FieldParams{"collect_type", 's', "collect type", ""},
		st.FieldParams{"collect_value", 'i', "collect value", "%"},
	})

	ds.Log().Debug("device.desiredState started (%s)", ds.Name())
	return nil
}
