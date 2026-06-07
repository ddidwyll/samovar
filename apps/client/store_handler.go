package client

import (
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"net/http"
)

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
	case "/store/fields":
		return sendJson(rw, stateFields.ToJson())
	case "/telemetry/scheme":
		return sendJson(rw, inter.TelemetryScheme(sh))
	default:
		return sendJson(rw, []byte(`{"ack":"ok"}`))
	}
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
