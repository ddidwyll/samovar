package calc

import (
	"samovar/lib/change"
	"samovar/lib/val"

	"strings"
)

type report = change.Report
type value = val.Val

type fieldId struct {
	stateKey string
	fieldKey string
}

type args map[fieldId]value
type calc func(args) value

type watcher struct {
	fields []fieldId
	target fieldId
	calc   calc
}

type State struct {
	watchers []watcher
}

func New(watchers ...watcher) State {
  return State{watchers: watchers}
}

func (s *State)

func fid(str string) fieldId {
	pair := strings.Split(str, ".")
	return fieldId{pair[0], pair[1]}
}
