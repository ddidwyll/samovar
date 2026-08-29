package telemetry

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
	"strings"
)

type Actor interface {
	Name() gen.Atom
	Call(any, any) (any, error)
	Send(any, any) error
}

type actor struct {
	id   gen.Atom
	name string
}

type relation struct {
	from   gen.Atom
	action string
	to     gen.Atom
}

type relations map[relation]int64

type actors map[gen.Atom]string

type state struct {
	act.Actor
	actors    actors
	relations relations
}

func newState() gen.ProcessBehavior {
	return &state{}
}

func (s *state) Init(_ ...any) error {
	s.relations = make(relations)
	s.actors = make(actors)
	return nil
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	switch r := msg.(type) {
	case relation:
		_, toExists := s.actors[r.to]
		_, fromExists := s.actors[r.from]

		if !toExists || !fromExists {
			err := fmt.Sprintf("Invalid telemetry actor relation: %s>%s>%s", r.from, r.action, r.to)
			return errors.New(err)
		} else {
			if current, has := s.relations[r]; has {
				s.relations[r] = current + 1
			} else {
				s.relations[r] = 1
			}
		}
	case actor:
		s.actors[r.id] = r.name
	default:
		return errors.New("Unexpected telemetry message")
	}
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, _ any) (any, error) {
	lines := make([]string, 0, len(s.relations))

	for relation, count := range s.relations {
		to := relation.to
		from := relation.from
		toName, toExists := s.actors[to]
		fromName, fromExists := s.actors[from]
		if !toExists || !fromExists {
			return lines, errors.New("Unexpected telemetry actor")
		}
		arrow := "-->"
		if strings.HasPrefix(relation.action, "request") {
			arrow = "<-->"
		}
		line := fmt.Sprintf(
			"%s%s %s|%s #%d| %s%s",
			string(from),
			fromName,
			arrow,
			relation.action,
			count,
			string(to),
			toName,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

func RegisterActor(a Actor, name string) {
	if err := a.Send(gen.Atom("telemetry_state"), actor{a.Name(), name}); err != nil {
		panic(err)
	}
}

func LogRelation(a Actor, action string, from, to gen.Atom) error {
	return a.Send(gen.Atom("telemetry_state"), relation{from, action, to})
}

func BuildScheme(a Actor) []byte {
	scheme, _ := a.Call(gen.Atom("telemetry_state"), nil)
	schemeStr, _ := scheme.(string)
	return []byte(schemeStr)
}
