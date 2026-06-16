package client

import (
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"encoding/json"
	"fmt"
	"net/http"
)

type errMsg struct {
	message string
	ack     string
}

func newStoreHandler() gen.ProcessBehavior {
	return &storeHandler{}
}

type storeHandler struct{ act.WebWorker }

type resp = http.ResponseWriter
type req = http.Request

func (sh *storeHandler) Init(_ ...any) error {
	inter.RegisterActor(sh, "[client.store.handler]")
	sh.Log().Debug("client.store_handler started")
	return nil
}

func (sh *storeHandler) HandleGet(_ gen.PID, rw resp, r *req) error {
	switch r.RequestURI {
	case "/store":
		fmt.Printf("Request: %#v", *r)
		return sendJson(rw, []byte(`{"ack":"error"}`))
	case "/store/fields":
		return sendJson(rw, stateFields.ToJson())
	case "/telemetry/scheme":
		return sendText(rw, inter.TelemetryScheme(sh))
	default:
		return sendJson(rw, []byte(`{"ack":"error"}`))
	}
}

func (sh *storeHandler) HandlePatch(_ gen.PID, rw resp, r *req) error {
	var patch map[string]string

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		return sendJsonError(rw, err)
	}
	defer r.Body.Close()

	if res, err := inter.Call(sh, patch, "patch", "client_desired_state"); err != nil {
		fmt.Printf("store.HandlePatch.Call: %#v, %#v\n", err, res)

		return sendJson(rw, []byte(`{"ack":"error"}`))
	} else {
		fmt.Printf("store.HandlePatch.Call: %#v, %#v\n", err, res)

		return sendJson(rw, []byte(`{"ack":"ok"}`))
	}
}

func sendJsonError(rw http.ResponseWriter, err error) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)
	errStr := fmt.Sprintf("Server error: %s", err)
	json, _ := json.Marshal(errMsg{errStr, "error"})
	_, err = rw.Write(json)

	return err
}

func sendJson(rw http.ResponseWriter, data []byte) (err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)

	return err
}

func sendText(rw resp, data []byte) (err error) {
	_, err = rw.Write(data)

	return err
}

// func (sh *storeHandler) HandleMessage(_ gen.PID, msg any) error {
//   sh.Log().Info("client.storeHandler received message: %#v", msg)
//   return nil
// }
