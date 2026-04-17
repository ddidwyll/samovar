package client

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"bytes"
	"encoding/json"
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
	info, _ := sh.Info()
	return sendJson(rw, info)
}

func sendJson(rw resp, data any) (err error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if err = enc.Encode(data); err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(buf.Bytes())

	return err
}

// func (sh *storeHandler) HandleMessage(_ gen.PID, msg any) error {
//   sh.Log().Info("client.storeHandler received message: %#v", msg)
//   return nil
// }
