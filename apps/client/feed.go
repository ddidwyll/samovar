package client

import (
	"samovar/lib/change"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/meta/sse"

	"errors"
	"fmt"
)

type conns map[gen.Alias]bool

type feed struct {
	act.Actor
	connections conns
	counter     int64
}

func newFeed() gen.ProcessBehavior {
	return &feed{
		connections: make(conns),
	}
}

func (f *feed) Init(_ ...any) (err error) {
	// f.SendAfter(f.PID(), "tick", 2*time.Second)
	f.Log().Debug("client.feed started (%s)", f.Name())
	return
}

func (f *feed) HandleMessage(from gen.PID, msg any) error {
	switch m := msg.(type) {
	case sse.MessageConnect:
		f.Log().Debug("New SSE connection: %s (remote: %s)", m.ID, m.RemoteAddr)
		f.connections[m.ID] = true

		connLen := len(f.connections)
		welcome := sse.Message{
			Event: "welcome",
			Data:  []byte(fmt.Sprintf("Connected! You are client #%d", connLen)),
			MsgID: fmt.Sprintf("%d", f.counter),
		}
		f.SendAlias(m.ID, welcome)

	case sse.MessageDisconnect:
		f.Log().Debug("SSE disconnected %s", m.ID)
		delete(f.connections, m.ID)

	case sse.MessageLastEventID:
		f.Log().Debug("Client reconnected with Last-Event-ID: %s", m.LastEventID)

	case change.Report:
		if m.LastFrom() == "client_state" {
			return f.broadcastChanges(m)
		}

	// case string:
	// 	if m == "tick" {
	// 		now := time.Now().Format("15.04.05")
	// 		conns := f.connections

	// 		for connID := range conns {
	// 			msg := sse.Message{
	// 				Event: "time",
	// 				Data:  []byte(fmt.Sprintf(`{"counter": %d, "time": "%s", "clients": %d}`, counter, now, len(conns))),
	// 				MsgID: fmt.Sprintf("%d", counter),
	// 			}
	// 			if err := f.SendAlias(connID, msg); err != nil {
	// 				f.Log().Error("Failed to send to %s: %s", connID, err)
	// 			}
	// 		}

	// 		f.SendAfter(f.PID(), "tick", 2*time.Second)
	// 	} else {
	// 		f.Log().Warning("client.feed received unexpected message: %s", m)
	// 	}
	default:
		err := fmt.Sprintf("client.feed: unexpected message (%#v)", msg)
		return errors.New(err)
	}

	return nil
}

func (f *feed) broadcastChanges(r change.Report) error {
	msg := sse.Message{
		Event: "change",
		Data:  r.ToJson(),
		MsgID: fmt.Sprintf("%d", f.counter),
	}

	for connID := range f.connections {
		if err := f.SendAlias(connID, msg); err != nil {
			return err
		}
	}

	f.counter += 1
	f.Log().Info("client.feed[%s]: %s", r.FormatField(), r.FormatValue())
	return nil
}
