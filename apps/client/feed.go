package client

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/meta/sse"

	"fmt"
	"time"
)

type conns map[gen.Alias]bool

type feed struct {
	act.Actor
	connections conns
	counter     uint
}

func newFeed() gen.ProcessBehavior {
	return &feed{
		connections: make(conns),
	}
}

func (f *feed) Init(_ ...any) (err error) {
	f.SendAfter(f.PID(), "tick", 2*time.Second)
	f.Log().Debug("client.feed started (%s)", f.Name())
	return
}

func (f *feed) HandleMessage(from gen.PID, message any) error {
	switch m := message.(type) {
	case sse.MessageConnect:
		f.Log().Debug("New SSE connection: %s (remote: %s)", m.ID, m.RemoteAddr)
		f.connections[m.ID] = true

		connLen := len(f.connections)
		welcome := sse.Message{
			Event: "welcome",
			Data:  []byte(fmt.Sprintf("Connected! You are client #%d", connLen)),
			MsgID: "0",
		}
		f.SendAlias(m.ID, welcome)

	case sse.MessageDisconnect:
		f.Log().Debug("SSE disconnected %s", m.ID)
		delete(f.connections, m.ID)

	case sse.MessageLastEventID:
		f.Log().Debug("Client reconnected with Last-Event-ID: %s", m.LastEventID)

	case string:
		if m == "tick" {
			f.counter++
			now := time.Now().Format("15.04.05")
			conns := f.connections
			counter := f.counter

			for connID := range conns {
				msg := sse.Message{
					Event: "time",
					Data:  []byte(fmt.Sprintf(`{"counter": %d, "time": "%s", "clients": %d}`, counter, now, len(conns))),
					MsgID: fmt.Sprintf("%d", counter),
				}
				if err := f.SendAlias(connID, msg); err != nil {
					f.Log().Error("Failed to send to %s: %s", connID, err)
				}
			}

			f.SendAfter(f.PID(), "tick", 2*time.Second)
		} else {
			f.Log().Warning("client.feed received unexpected message: %s", m)
		}
	}

	return nil
}
