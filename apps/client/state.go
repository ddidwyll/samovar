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
	st.DefFieldUnit("t_top", 'f', "temp top", "°C"),
	st.DefFieldUnit("t_mid", 'f', "temp middle", "°C"),
	st.DefFieldUnit("t_btm", 'f', "temp bottom", "°C"),
	st.DefFieldUnit("mid_stable_temp", 'f', "middle stable temp", "°C"),
	st.DefFieldUnit("mid_stable_diff", 'f', "middle first stable temp diff", "°C"),
	st.DefField("is_mid_stable", 's', "is middle temp stable"),
	st.DefFieldUnit("power", 'i', "power", "W"),
	st.DefFieldUnit("net_power", 'i', "net power", "W"),
	st.DefFieldUnit("power_diff", 'i', "power diff", "%"),
	st.DefFieldUnit("collect", 's', "collect", "%"),
	st.DefFieldUnit("collection_speed", 's', "collection speed", "g/h"),
	st.DefField("reflux_ratio", 'f', "reflux_ratio"),
	st.DefField("collect_synced", 's', "is collect synced"),
	st.DefField("body_collected", 's', "collected body"),
	st.DefField("head_collected", 's', "collected head"),
	st.DefField("recyc_collected", 's', "collected recyc"),
	st.DefField("waste_collected", 's', "collected waste"),
	st.DefFieldUnit("body_average_speed", 'i', "body average speed", "g/h"),
	st.DefFieldUnit("head_average_speed", 'i', "head average speed", "g/h"),
	st.DefFieldUnit("recyc_average_speed", 'i', "recyc average speed", "g/h"),
	st.DefFieldUnit("press", 'f', "press", "mm"),
	st.DefField("devices", 'a', "devices"),
	st.DefField("device_id", 's', "device id"),
	st.DefField("scripts", 'a', "scripts"),
	st.DefField("script_mode", 's', "script mode"),
	st.DefField("error", 's', "error"),
	st.DefField("ready", 's', "ready"),
	// st.DefField("last_mqtt_tx", 's', "last mqtt tx"),
	// st.DefField("last_mqtt_rx", 's', "last mqtt rx"),
	st.DefField("last_mqtt_ping", 's', "last mqtt ping"),
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
