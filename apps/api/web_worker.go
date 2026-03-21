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
	return &webWorker{}
}

type webWorker struct{ act.WebWorker }

// Init invoked on a start this process.
func (ww *webWorker) Init(args ...any) error {
	ww.SendAfter(ww.PID(), "tick", 2*time.Second)
	ww.Log().Debug("started web worker process with args %v", args)
	return nil
}

// Handle GET requests. For the other HTTP methods (POST, PATCH, etc)
// you need to add the accoring callback-method implementation. See act.WebWorkerBehavior.

func (ww *webWorker) HandleGet(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	var buf bytes.Buffer

	ww.Log().Debug("got HTTP request %q", request.URL.Path)
	writer.Header().Set("Content-Type", "application/json")
	// response JSON message with information about this process
	info, _ := ww.Info()
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(info)
	writer.Write(buf.Bytes())
	return nil
}

func (ww *webWorker) sseStateAddConn(alias gen.Alias) {
	if err := ww.Send(sseStateProcessName, sseStateReqAddConn{alias}); err != nil {
		ww.Log().Error("%s", err)
		panic(err)
	}
}

func (ww *webWorker) sseStateDelConn(alias gen.Alias) {
	if err := ww.Send(sseStateProcessName, sseStateReqDelConn{alias}); err != nil {
		ww.Log().Error("%s", err)
		panic(err)
	}
}

func (ww *webWorker) sseStateCountIncr() {
	if err := ww.Send(sseStateProcessName, sseStateReqCountIncr{}); err != nil {
		ww.Log().Error("%s", err)
		panic(err)
	}
}

func (ww *webWorker) sseStateGet(key string) any {
	result, err := ww.Call(sseStateProcessName, sseStateReqGet{key})

	if err == nil {
		return result
	} else {
		ww.Log().Error("%s", err)
		panic(err)
	}
}

func (ww *webWorker) sseStateGetConns() sseStateConns {
	if m, ok := ww.sseStateGet("connections").(sseStateConns); ok {
		return m
	} else {
		panic("Unexpected sseState state value")
	}
}

func (ww *webWorker) sseStateGetCounter() uint {
	if c, ok := ww.sseStateGet("counter").(uint); ok {
		return c
	} else {
		panic("Unexpected sseState state value")
	}
}

func (ww *webWorker) sseStateGetConnLen() int {
	return len(ww.sseStateGetConns())
}

func (ww *webWorker) HandleMessage(from gen.PID, message any) error {
	switch m := message.(type) {
	case sse.MessageConnect:
		ww.Log().Debug("New SSE connection: %s (remote: %s)", m.ID, m.RemoteAddr)
		ww.sseStateAddConn(m.ID)

		connLen := ww.sseStateGetConnLen()
		welcome := sse.Message{
			Event: "welcome",
			Data:  []byte(fmt.Sprintf("Connected! You are client #%d", connLen)),
			MsgID: "0",
		}
		ww.SendAlias(m.ID, welcome)

	case sse.MessageDisconnect:
		ww.Log().Debug("SSE disconnected %s", m.ID)
		ww.sseStateDelConn(m.ID)

	case sse.MessageLastEventID:
		ww.Log().Debug("Client reconnected with Last-Event-ID: %s", m.LastEventID)

	case string:
		if m == "tick" {
			ww.sseStateCountIncr()
			now := time.Now().Format("15.04.05")
			conns := ww.sseStateGetConns()
			counter := ww.sseStateGetCounter()

			for connID := range conns {
				msg := sse.Message{
					Event: "time",
					Data:  []byte(fmt.Sprintf(`{"counter": %d, "time": "%s", "clients": %d}`, counter, now, len(conns))),
					MsgID: fmt.Sprintf("%d", counter),
				}
				if err := ww.SendAlias(connID, msg); err != nil {
					ww.Log().Error("Failed to send to %s: %s", connID, err)
				}
			}

			ww.SendAfter(ww.PID(), "tick", 2*time.Second)
		} else {
			ww.Log().Warning("Unexpected webWorker message: %s", m)
		}
	}

	return nil
}
