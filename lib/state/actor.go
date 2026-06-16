package state

import (
	"samovar/lib/change"
	"samovar/lib/inter"

	"ergo.services/ergo/act"

	"errors"
	"fmt"
	"time"
)

type StateActor struct {
	act.Actor
	data *State
}

func (sa *StateActor) InitState(name string, fields Fields) {
	inter.RegisterActor(sa, name)
	sa.data = BuildState(fields)
}

func (sa *StateActor) BulkChange(reqs []request) ([]report, error) {
	reports := make([]report, 0, len(reqs))
	lastFrom := string(sa.Name())

	for _, req := range reqs {
		if rep, err := sa.data.Change(req); err != nil {
			return reports, err
		} else {
			if rep.Changed {
				reports = append(reports, rep.AddFrom(lastFrom))
			}
		}
	}

	return reports, nil
}

func (sa *StateActor) BulkChangeFromMap(kv map[string]string, from, producer string) error {
	ts := time.Now().UnixMicro()
	reqs := make([]request, 0, len(kv))
	for k, v := range kv {
		req := change.NewRequest(k, v, from, ts)
		reqs = append(reqs, req)
	}
	return sa.HandleChangeRequests(reqs, producer)
}

func (sa *StateActor) HandleChangeRequests(req any, producer string) error {
	switch r := req.(type) {
	case request:
		if report, err := sa.data.Change(r); err != nil {
			return err
		} else {
			lastFrom := string(sa.Name())
			return inter.Send(sa, report.AddFrom(lastFrom), "report", producer)
		}
	case []request:
		if reports, err := sa.BulkChange(r); err != nil {
			return err
		} else {
			return inter.Send(sa, reports, "reports", producer)
		}
	default:
		err := fmt.Sprintf("Unexpected state change request: %#v", req)
		return errors.New(err)
	}
}

func (sa *StateActor) HandleDataRequest(req any) (any, error) {
	switch k := req.(type) {
	case []string:
		return sa.data.KeyVals(k...)
	case string:
		return sa.data.FetchField(k)
	case nil:
		return sa.data.Entries(), nil
	default:
		return nil, errors.New("Unexpected state data request")
	}
}
