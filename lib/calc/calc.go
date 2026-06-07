package calc

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/state"
	"samovar/lib/val"

	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type report = change.Report
type Value = val.Val

type fieldId struct {
	stateKey string
	fieldKey string
}

type changes map[string]map[string]Value

type Args = state.KeyVals
type Results = Args

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

	for stateKey, changes := range changes {
		requests := make([]change.Request, 0, len(changes))

		for fieldKey, value := range changes {
			from := string(c.process.Name())
			request := r.NewRequest(from, fieldKey, value)
			requests = append(requests, request)
		}

		if err = inter.Send(c.process, requests, "requests", stateKey); err != nil {
			return err
		}
	}

	return nil
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
			if value.IsNil() {
				continue
			}
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
