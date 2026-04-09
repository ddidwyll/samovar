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

func (ah *apiHandler) Init(args ...any) error {
	ah.Log().Debug("started api handler process with args %v", args)
	return nil
}

func (ah *apiHandler) HandleGet(from gen.PID, w http.ResponseWriter, r *http.Request) error {
	var buf bytes.Buffer

	ah.Log().Debug("got HTTP r %q", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	info, _ := ah.Info()
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(info)
	w.Write(buf.Bytes())
	return nil
}
