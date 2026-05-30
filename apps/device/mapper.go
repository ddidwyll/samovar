package device

import (
	// "samovar/lib/val"
	"samovar/lib/change"
)

type request = change.Request
type mapper func(request) (request, bool)

var mappers = []mapper{tryAsIs}

var asIsMap = map[string]string{
	"term_d":  "t_top",
	"term_c":  "t_mid",
	"term_k":  "t_btm",
	"power_m": "power",
	"otbor":   "collect",
}

func mapFromRaw(req request) (request, bool) {
	for _, fn := range mappers {
		if req, mapped := fn(req); mapped {
			return req, true
		}
	}
	return req, false
}

func tryAsIs(req request) (request, bool) {
	if newKey, has := asIsMap[req.Key]; has {
		req.Key = newKey
		return req, true
	}
	return req, false
}
