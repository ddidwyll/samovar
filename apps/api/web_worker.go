package api

import (
	"bytes"
	"encoding/json"
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/meta/sse"
	"fmt"
	"net/http"
	"time"
)

func createWebWorker() gen.ProcessBehavior {
	return &WebWorker{}
}

type WebWorker struct {
	act.WebWorker
	connections map[gen.Alias]bool
	counter     uint
}

// Init invoked on a start this process.
func (w *WebWorker) Init(args ...any) error {
	w.connections = make(map[gen.Alias]bool)
	w.SendAfter(w.PID(), "tick", 2*time.Second)
	w.Log().Info("started web worker process with args %v", args)
	return nil
}

// Handle GET requests. For the other HTTP methods (POST, PATCH, etc)
// you need to add the accoring callback-method implementation. See act.WebWorkerBehavior.

func (w *WebWorker) HandleGet(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	var buf bytes.Buffer

	w.Log().Info("got HTTP request %q", request.URL.Path)
	writer.Header().Set("Content-Type", "application/json")
	// response JSON message with information about this process
	info, _ := w.Info()
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(info)
	writer.Write(buf.Bytes())
	return nil
}

func (w *WebWorker) HandleMessage(from gen.PID, message any) error {
	switch m := message.(type) {
	case sse.MessageConnect:
		w.Log().Info("New SSE connection: %s (remote: %s)", m.ID, m.RemoteAddr)
		w.connections[m.ID] = true

		welcome := sse.Message{
			Event: "welcome",
			Data:  []byte(fmt.Sprintf("Connected! You are client #%d", len(w.connections))),
			MsgID: "0",
		}
		w.SendAlias(m.ID, welcome)

	case sse.MessageDisconnect:
		w.Log().Info("SSE disconnected %s", m.ID)
		delete(w.connections, m.ID)

	case sse.MessageLastEventID:
		w.Log().Info("Client reconnected with Last-Event-ID: %s", m.LastEventID)

	case string:
		if m == "tick" {
			w.counter++
			now := time.Now().Format("15.04.05")

			for connID := range w.connections {
				msg := sse.Message{
					Event: "time",
					Data:  []byte(fmt.Sprintf(`{"counter": %d, "time": "%s", "clients": %d}`, w.counter, now, len(w.connections))),
					MsgID: fmt.Sprintf("%d", w.counter),
				}
				if err := w.SendAlias(connID, msg); err != nil {
					w.Log().Error("Failed to send to %s: %s", connID, err)
				}
			}

			w.SendAfter(w.PID(), "tick", 2*time.Second)
		} else {
			w.Log().Warning("Unexpected WebWorker message: %s", m)
		}
	}

	return nil
}
