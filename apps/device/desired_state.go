package device

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type desiredState struct{ st.StateActor }

var desiredStateFields = st.Fields{
	st.FieldParams{"power", 'i', "power", "W"},
	// st.FieldParams{"collect_type", 's', "collect type", ""},
	// st.FieldParams{"collect_value", 'i', "collect value", "%"},
}

func newDesiredState() gen.ProcessBehavior {
	return &desiredState{}
}

func (ds *desiredState) Init(_ ...any) error {
	ds.InitState("[(device.desired_state)]", desiredStateFields)
	return nil
}

func (ds *desiredState) HandleMessage(_ gen.PID, msg any) error {
  ds.Log().Info("device.desiredState.HandleMessage.msg: %v", msg)
	return ds.HandleChangeRequests(msg, "device_producer")
}

func (ds *desiredState) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return ds.HandleDataRequest(req)
}
