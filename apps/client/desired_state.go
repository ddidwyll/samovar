package client

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type desiredState struct{ st.StateActor }

var desiredStateFields = st.Fields{
	st.FieldParams{"device_id", 's', "device id", ""},
	st.FieldParams{"script", 's', "script", ""},
}

func newDesiredState() gen.ProcessBehavior {
	return &desiredState{}
}

func (ds *desiredState) Init(_ ...any) error {
	ds.InitState("[(client.desired_state)]", desiredStateFields)
	return nil
}

func (ds *desiredState) HandleMessage(_ gen.PID, msg any) error {
	ds.Log().Debug("client.desiredState.HandleMessage.msg: %v", msg)
	return ds.HandleChangeRequests(msg, "client_producer")
}

func (ds *desiredState) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	if kv, ok := req.(map[string]string); ok {
		ds.Log().Debug("client.desiredState.HandleCall.req: %v", req)
		return "", ds.BulkChangeFromMap(kv, "client_store", "client_producer")
	} else {
		return ds.HandleDataRequest(req)
	}
}
