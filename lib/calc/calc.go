package calc

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/state"
	"samovar/lib/val"

	"ergo.services/ergo/act"

	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

type report = change.Report
type Value = any

type fieldId struct {
	stateKey string
	fieldKey string
}

type Args = state.KeyVals
type Results = map[string]Value

type changes map[string]Results
type fields map[string]map[string]bool

type ApplyFn func(string, Value)
type FetchFn func(string, string) val.Val

type fieldsCalcFn func(Args, ApplyFn)
type reportCalcFn func(report, FetchFn, ApplyFn)

type fieldsWatcher struct {
	fields fields
	calc   fieldsCalcFn
}

type reportWatcher struct {
	stateKey string
	fieldKey string
	calc     reportCalcFn
}

type CalcActor struct {
	act.Actor
	fieldsWatchers []fieldsWatcher
	reportWatchers []reportWatcher
	applyFn        ApplyFn
}

func (ca *CalcActor) InitCalc(name string) {
	inter.RegisterActor(ca, name)
	ca.reportWatchers = make([]reportWatcher, 0)
	ca.fieldsWatchers = make([]fieldsWatcher, 0)
}

func (ca *CalcActor) SetApplyFn(fn ApplyFn) {
	ca.applyFn = fn
}

func (ca *CalcActor) WatchFieldAs(field, asField string) {
	fn := func(args Args, apply ApplyFn) {
		for _, value := range args {
			apply(asField, value)
		}
	}
	ca.WatchFields(fn, []string{field}...)
}

func (ca *CalcActor) WatchFields(fn fieldsCalcFn, fieldStrings ...string) {
	fields := make(fields)

	for _, fieldIdStr := range fieldStrings {
		fid := fid(fieldIdStr)
		if fields[fid.stateKey] == nil {
			fields[fid.stateKey] = make(map[string]bool)
		}
		fields[fid.stateKey][fid.fieldKey] = true
	}

	watcher := fieldsWatcher{fields, fn}

	ca.fieldsWatchers = append(ca.fieldsWatchers, watcher)
}

func (ca *CalcActor) HandleChangeReports(msg any) error {
	r, ok := msg.(report)
	if !ok {
		err := fmt.Sprintf("Invalid change report: %v", msg)
		return errors.New(err)
	}

	changes, err := ca.performWatchers(r)
	if err != nil {
		return err
	}

	for stateKey, results := range changes {
		if err := ca.sendRequestsWithReport(stateKey, results, r); err != nil {
			return err
		}
	}

	return nil
}

func (ca *CalcActor) SendRequest(stateKey, fieldKey string, value Value) error {
	results := make(Results, 1)
	results[fieldKey] = value
	return ca.SendRequests(stateKey, results)
}

func (ca *CalcActor) SendRequests(stateKey string, results Results) error {
	ts := time.Now().UnixMicro()

	requests := make([]change.Request, 0, len(results))

	for fieldKey, value := range results {
		requests = append(requests, ca.BuildRequest(fieldKey, value, ts))
	}

	return inter.Send(ca, requests, "requests", stateKey)
}

func (ca *CalcActor) sendRequestsWithReport(stateKey string, results Results, r report) error {
	requests := make([]change.Request, 0, len(results))

	for fieldKey, value := range results {
		requests = append(requests, ca.buildRequestFromReport(fieldKey, value, r))
	}

	return inter.Send(ca, requests, "requests", stateKey)
}

func (ca *CalcActor) buildRequestFromReport(fieldKey string, value Value, r report) change.Request {
	from := string(ca.Name())
	return r.NewRequest(from, fieldKey, value)
}

func (ca *CalcActor) BuildRequest(fieldKey string, value Value, ts int64) change.Request {
	from := string(ca.Name())
	return change.NewRequest(fieldKey, value, from, ts)
}

func (ca *CalcActor) mustFetchField(stateKey, fieldKey string) val.Val {
	if v, err := inter.Call(ca, fieldKey, "field", stateKey); err == nil {
		return v.(val.Val)
	} else {
		panic(err)
	}
}

func (ca *CalcActor) performWatchers(r report) (changes, error) {
	allResults := make(changes)
	stateKey := r.LastFrom()

	compileResults := func(results Results) {
		for fieldIdStr, value := range results {
			target := fid(fieldIdStr)
			sk := target.stateKey
			fk := target.fieldKey
			if allResults[sk] == nil {
				allResults[sk] = make(Results)
			}
			allResults[sk][fk] = value
		}
	}

	withResults := func(cb func(ApplyFn)) {
		results := make(Results)

		applyFn := func(field string, value Value) {
			results[field] = value
		}

		cb(applyFn)

		compileResults(results)
	}

	for _, watcher := range ca.fieldsWatchers {
		if !watcher.match(stateKey, r.Key) {
			continue
		}

		args, err := ca.buildArgs(watcher)
		if err != nil {
			return allResults, err
		}

		if ca.applyFn != nil {
			watcher.calc(args, ca.applyFn)
		} else {
			results := make(Results)

			applyFn := func(field string, value Value) {
				results[field] = value
			}
			watcher.calc(args, applyFn)

			compileResults(results)
		}
	}

	for _, watcher := range ca.reportWatchers {
		if !watcher.match(stateKey, r.Key) {
			continue
		}

		if ca.applyFn != nil {
			watcher.calc(r, ca.mustFetchField, ca.applyFn)
		} else {
			withResults(func(applyFn ApplyFn) {
				watcher.calc(r, ca.mustFetchField, applyFn)
			})
		}
	}

	return allResults, nil
}

func (ca *CalcActor) buildArgs(w fieldsWatcher) (Args, error) {
	args := make(Args)

	for stateKey, fieldKeyMap := range w.fields {
		fieldKeys := slices.Collect(maps.Keys(fieldKeyMap))

		kv, err := inter.Call(ca, fieldKeys, "fields", stateKey)
		if err != nil {
			return args, err
		}

		stateArgs, ok := kv.(Args)
		if !ok {
			e := fmt.Sprintf("calc.buildArgs.stateArgs: %s, %#kv", ok, kv)
			return args, errors.New(e)
		}

		for key, val := range stateArgs {
			argKey := fidFrom(stateKey, key).String()
			args[argKey] = val
		}
	}
	return args, nil
}

func (w fieldsWatcher) match(stateKey, fieldKey string) bool {
	return w.fields[stateKey][fieldKey]
}

func (w reportWatcher) match(stateKey, fieldKey string) bool {
	return w.stateKey == stateKey && w.fieldKey == fieldKey
}

func fid(str string) fieldId {
	pair := strings.Split(str, ".")
	return fieldId{pair[0], pair[1]}
}

func (fid fieldId) String() string {
	return fmt.Sprintf("%s.%s", fid.stateKey, fid.fieldKey)
}

func fidFrom(stateKey, fieldKey string) fieldId {
	return fieldId{stateKey, fieldKey}
}
