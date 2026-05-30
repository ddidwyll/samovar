package calc

import (
	"samovar/lib/change"
	"samovar/lib/val"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type report = change.Report
type value = val.Val

type fieldId struct {
	stateKey string
	fieldKey string
}

type changes map[string]map[string]value

type Args map[string]value
type calcFn func(Args) (value, error)

type fields map[string]map[string]bool

type watcher struct {
	fieldIds []fieldId
	fields   fields
	target   fieldId
	calc     calcFn
}

type process interface {
	Call(to any, message any) (any, error)
	Send(to any, message any) error
	Name() gen.Atom
}

type Calc struct {
	watchers []watcher
	process  process
}

func NewCalc(p process) *Calc {
	return &Calc{make([]watcher, 0), p}
}

func (c *Calc) Watch(fn calcFn, target string, fieldStrings ...string) {
	fields := make(fields)
	fieldIds := make([]fieldId, 0, len(fieldStrings))

	for _, fieldIdStr := range fieldStrings {
		fid := fid(fieldIdStr)
		fields[fid.stateKey][fid.fieldKey] = true
		fieldIds = append(fieldIds, fid)
	}

	watcher := watcher{fieldIds, fields, fid(target), fn}

	c.watchers = append(c.watchers, watcher)
}

func (c *Calc) HandleReport(maybeReport any) error {
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

		if err = c.process.Send(gen.Atom(stateKey), requests); err != nil {
			return err
		}
	}

	return nil
}

func (c *Calc) performWatchers(r report) (results changes, err error) {
	stateKey := r.LastFrom()

	for _, watcher := range c.watchers {
		if !watcher.match(stateKey, r.Key) {
			continue
		}

		args, err := watcher.buildArgs(c.process)
		if err != nil {
			return results, err
		}

		if value, err := watcher.calc(args); err != nil {
			return results, err
		} else {
			results[watcher.target.stateKey][watcher.target.fieldKey] = value
		}
	}

	return results, nil
}

func (w watcher) match(stateKey, fieldKey string) bool {
	return w.fields[stateKey][fieldKey]
}

func (w watcher) buildArgs(p process) (result Args, err error) {
	for stateKey, fieldKeyMap := range w.fields {
		fieldKeys := slices.Collect(maps.Keys(fieldKeyMap))

		kv, err := p.Call(gen.Atom(stateKey), fieldKeys)
		if err != nil {
			return result, err
		}

		stateArgs, ok := kv.(Args)
		if !ok {
			return result, errors.New("Failed to fetch calc args")
		}

		for key, val := range stateArgs {
			argKey := fidFrom(stateKey, key).String()
			result[argKey] = val
		}
	}
	return result, nil
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
