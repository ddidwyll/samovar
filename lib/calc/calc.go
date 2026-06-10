package calc

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/state"
	// "samovar/lib/val"

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

type ApplyFn func(string, Value)
type calcFn func(Args, ApplyFn)

type fields map[string]map[string]bool

type watcher struct {
	fields fields
	calc   calcFn
}

type process = inter.Actor

type Config struct {
	watchers []watcher
	process  process
}

func NewConfig(p process, name string) *Config {
	inter.RegisterActor(p, name)
	return &Config{make([]watcher, 0), p}
}

func (c *Config) WatchAs(field, asField string) {
	fn := func(args Args, apply ApplyFn) {
		for _, value := range args {
			apply(asField, value)
		}
	}
	c.Watch(fn, []string{field}...)
}

func (c *Config) Watch(fn calcFn, fieldStrings ...string) {
	fields := make(fields)

	for _, fieldIdStr := range fieldStrings {
		fid := fid(fieldIdStr)
		if fields[fid.stateKey] == nil {
			fields[fid.stateKey] = make(map[string]bool)
		}
		fields[fid.stateKey][fid.fieldKey] = true
	}

	watcher := watcher{fields, fn}

	c.watchers = append(c.watchers, watcher)
}

func (c *Config) HandleReport(maybeReport any) error {
	r, ok := maybeReport.(report)
	if !ok {
		err := fmt.Sprintf("Invalid change report: %v", maybeReport)
		return errors.New(err)
	}

	changes, err := c.performWatchers(r)
	if err != nil {
		return err
	}

	for stateKey, results := range changes {
		if err := c.sendRequestsWithReport(stateKey, results, r); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) SendRequest(stateKey, fieldKey string, value Value) error {
	results := make(Results, 1)
	results[fieldKey] = value
	return c.SendRequests(stateKey, results)
}

func (c *Config) SendRequests(stateKey string, results Results) error {
	ts := time.Now().UnixMicro()

	requests := make([]change.Request, 0, len(results))

	for fieldKey, value := range results {
		requests = append(requests, c.BuildRequest(fieldKey, value, ts))
	}

	return inter.Send(c.process, requests, "requests", stateKey)
}

func (c *Config) sendRequestsWithReport(stateKey string, results Results, r report) error {
	requests := make([]change.Request, 0, len(results))

	for fieldKey, value := range results {
		requests = append(requests, c.buildRequestFromReport(fieldKey, value, r))
	}

	return inter.Send(c.process, requests, "requests", stateKey)
}

func (c *Config) buildRequestFromReport(fieldKey string, value Value, r report) change.Request {
	from := string(c.process.Name())
	return r.NewRequest(from, fieldKey, value)
}

func (c *Config) BuildRequest(fieldKey string, value Value, ts int64) change.Request {
	from := string(c.process.Name())
	return change.NewRequest(fieldKey, value, from, ts)
}

func (c *Config) performWatchers(r report) (changes, error) {
	allResults := make(changes)
	stateKey := r.LastFrom()

	for _, watcher := range c.watchers {
		if !watcher.match(stateKey, r.Key) {
			continue
		}

		args, err := watcher.buildArgs(c.process)
		if err != nil {
			return allResults, err
		}

		results := make(Results)
		applyFn := func(field string, value Value) {
			results[field] = value
		}
		watcher.calc(args, applyFn)

		for fieldIdStr, value := range results {
			// if value.IsNil() {
			// 	continue
			// }
			target := fid(fieldIdStr)
			sk := target.stateKey
			fk := target.fieldKey
			if allResults[sk] == nil {
				allResults[sk] = make(Results)
			}
			allResults[sk][fk] = value
		}
	}

	return allResults, nil
}

func (w watcher) match(stateKey, fieldKey string) bool {
	return w.fields[stateKey][fieldKey]
}

func (w watcher) buildArgs(p process) (Args, error) {
	args := make(Args)

	for stateKey, fieldKeyMap := range w.fields {
		fieldKeys := slices.Collect(maps.Keys(fieldKeyMap))

		kv, err := inter.Call(p, fieldKeys, "fields", stateKey)
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
