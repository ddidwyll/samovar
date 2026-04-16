package client

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"bytes"
	"encoding/json"
	"net/http"
)

func newApiHandler() gen.ProcessBehavior {
	return &apiHandler{}
}

type apiHandler struct{ act.WebWorker }

type resp = http.ResponseWriter
type req = http.Request

func (ah *apiHandler) Init(_ ...any) error {
	ah.Log().Debug("client.api_handler started")
	return nil
}

func (ah *apiHandler) HandleGet(_ gen.PID, rw resp, r *req) error {
	info, _ := ah.Info()
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
	_, err  = rw.Write(buf.Bytes())

	return err
}
