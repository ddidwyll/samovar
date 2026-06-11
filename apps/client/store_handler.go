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

	fmt.Printf("PATCH: %#v", patch)

	return sendJson(rw, []byte(`{"ack":"error"}`))
}

func sendJsonError(rw resp, err error) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)
	errStr := fmt.Sprintf("Server error: %s", err)
	json, _ := json.Marshal(errMsg{errStr, "error"})
	_, err = rw.Write(json)

	return err
}

func sendJson(rw resp, data []byte) (err error) {
	rw.Header().Set("Content-Type", "application/json")
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
