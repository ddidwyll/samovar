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
type calcFn func(Args) (Value, error)

type fields map[string]map[string]bool

type watcher struct {
	fields fields
	target fieldId
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

func (c *Config) Watch(fn calcFn, target string, fieldStrings ...string) {
	fields := make(fields)

	for _, fieldIdStr := range fieldStrings {
		fid := fid(fieldIdStr)
		if fields[fid.stateKey] == nil {
			fields[fid.stateKey] = make(map[string]bool)
		}
		fields[fid.stateKey][fid.fieldKey] = true
	}

	watcher := watcher{fields, fid(target), fn}

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

func Error(err string) (Value, error) {
	return val.Nil{}, errors.New(err)
}

func Skip() (Value, error) {
	return val.Nil{}, nil
}

func TryAsIs(args Args) (Value, error) {
	if len(args) != 1 {
		return Error("Invalid calc arguments")
	}
	for _, val := range args {
		return val, nil
	}
	return Skip()
}

func (c *Config) performWatchers(r report) (changes, error) {
	results := make(changes)
	stateKey := r.LastFrom()

	for _, watcher := range c.watchers {
		if !watcher.match(stateKey, r.Key) {
			continue
		}

		stateKey := watcher.target.stateKey
		fieldKey := watcher.target.fieldKey

		if results[stateKey] == nil {
			results[stateKey] = make(map[string]Value)
		}

		args, err := watcher.buildArgs(c.process)
		if err != nil {
			return results, err
		}

		if value, err := watcher.calc(args); err != nil {
			return results, err
		} else {
			if !value.IsNil() {
				results[stateKey][fieldKey] = value
			}
		}
	}

	return results, nil
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
