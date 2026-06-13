package client

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type desiredState struct{ st.StateActor }

var desiredStateFields = st.Fields{
	st.FieldParams{"device_id", 's', "device id", ""},
}

func newDesiredState() gen.ProcessBehavior {
	return &desiredState{}
}

func (ds *desiredState) Init(_ ...any) error {
	ds.InitState("[(client.desired_state)]", desiredStateFields)
	return nil
}

func (ds *desiredState) HandleMessage(_ gen.PID, req any) error {
	return ds.HandleChangeRequests(req, "client_producer")
}

func (ds *desiredState) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	if kv, ok := req.(map[string]string); ok {
		return "", ds.BulkChangeFromMap(kv, "client_store_patch", "client_producer")
	} else {
		return ds.HandleDataRequest(req)
	}
}
