package client

import (
	"samovar/lib/inter"
	st "samovar/lib/state"

	"ergo.services/ergo/gen"

	"errors"
)

type state struct{ st.StateActor }

func newState() gen.ProcessBehavior {
	return &state{}
}

var stateFields = st.Fields{
	st.FieldParams{"t_top", 'f', "t top", "°C"},
	st.FieldParams{"t_mid", 'f', "t middle", "°C"},
	st.FieldParams{"t_btm", 'f', "t bottom", "°C"},
	st.FieldParams{"mid_stable_temp", 'f', "middle stable temp", "°C"},
	st.FieldParams{"is_mid_stable", 's', "is middle temp stable", ""},
	st.FieldParams{"power", 'i', "power", "W"},
	st.FieldParams{"net_power", 'i', "net power", "W"},
	st.FieldParams{"power_diff", 'i', "power diff", "%"},
	st.FieldParams{"collect", 's', "collect", "%"},
	st.FieldParams{"collection_speed", 's', "collection speed", "g/h"},
	st.FieldParams{"reflux_ratio", 'f', "reflux_ratio", ""},
	st.FieldParams{"collect_synced", 's', "is collect synced", ""},
	st.FieldParams{"body_collected", 's', "collected body", ""},
	st.FieldParams{"head_collected", 's', "collected head", ""},
	st.FieldParams{"recyc_collected", 's', "collected recyc", ""},
	st.FieldParams{"body_average_speed", 'i', "body average speed", "g/h"},
	st.FieldParams{"head_average_speed", 'i', "head average speed", "g/h"},
	st.FieldParams{"recyc_average_speed", 'i', "recyc average speed", "g/h"},
	st.FieldParams{"press", 'f', "press", "mm"},
	st.FieldParams{"devices", 'a', "devices", ""},
	st.FieldParams{"device_id", 's', "device id", ""},
	st.FieldParams{"scripts", 'a', "scripts", ""},
	st.FieldParams{"script_mode", 's', "script mode", ""},
	st.FieldParams{"error", 's', "error", ""},
	st.FieldParams{"ready", 's', "ready", ""},
	// st.FieldParams{"last_mqtt_tx", 's', "last mqtt tx", ""},
	// st.FieldParams{"last_mqtt_rx", 's', "last mqtt rx", ""},
	st.FieldParams{"last_mqtt_ping", 's', "last mqtt ping", ""},
}

func (s *state) Init(_ ...any) error {
	s.InitState("[(client.state)]", stateFields)
	return nil
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("client.state.HandleMessage.msg: %s", msg)
	return s.HandleChangeRequests(msg, "client_producer")
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	s.Log().Debug("client.state.HandleCall.req: %s", req)
	return s.HandleDataRequest(req)
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("client.state.Terminate.reason: %s", reason)
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
