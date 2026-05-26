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
	f.Log().Debug("client.feed started (%s)", f.Name())
	return
}

func (f *feed) HandleMessage(from gen.PID, msg any) error {
	switch m := msg.(type) {
	case sse.MessageConnect:
		return f.handleConnect(m)

	case sse.MessageDisconnect:
		f.Log().Debug("SSE disconnected %s", m.ID)
		delete(f.connections, m.ID)

	case sse.MessageLastEventID:
		f.Log().Debug("Client reconnected with Last-Event-ID: %s", m.LastEventID)

	case change.Report:
		if m.LastFrom() == "client_state" {
			return f.broadcastFeed(m)
		}

	default:
		err := fmt.Sprintf("client.feed: unexpected message (%#v)", msg)
		return errors.New(err)
	}

	return nil
}

func (f *feed) handleConnect(m sse.MessageConnect) error {
	f.Log().Debug("New SSE connection: %s (remote: %s)", m.ID, m.RemoteAddr)

	if entries, err := stateEntries(f); err != nil {
		f.Log().Error("%v", err)
		return err
	} else {
		f.connections[m.ID] = true

		initMsg := sse.Message{
			Event: "init",
			Data:  entries.ToJson(),
			MsgID: fmt.Sprintf("%d", f.counter),
		}

		return f.SendAlias(m.ID, initMsg)
	}
}

func (f *feed) broadcastFeed(r change.Report) error {
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
	f.Log().Info("client.feed[%s]: %s", r.FieldName, r.FormatValue())
	return nil
}
