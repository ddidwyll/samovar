package session

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type desiredState struct{ st.StateActor }

var desiredStateFields = st.Fields{
	st.DefFieldUnit("power", 'i', "power", "W"),
	// st.DefField("collect_type", 's', "collect type"),
	// st.DefFieldUnit("collect_value", 'i', "collect value", "%"),
}

func newDesiredState() gen.ProcessBehavior {
	return &desiredState{}
}

func (ds *desiredState) Init(_ ...any) error {
	ds.InitState("[(session.desired_state)]", desiredStateFields)
	return nil
}

func (ds *desiredState) HandleMessage(_ gen.PID, msg any) error {
	ds.Log().Debug("session.desiredState.HandleMessage.msg: %v", msg)
	return ds.HandleChangeRequests(msg, "session_producer")
}

func (ds *desiredState) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return ds.HandleDataRequest(req)
}
