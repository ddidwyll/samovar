package inter

import (
	"samovar/apps/telemetry"

	"ergo.services/ergo/gen"

	"fmt"
	"strings"
)

type Actor = telemetry.Actor

func Send(a Actor, msg any, what, to string) error {
	toAtom := gen.Atom(to)

	if err := a.Send(toAtom, msg); err != nil {
		return logError(err, a.Name(), to, what)
	} else {
		action := fmt.Sprintf("send %s", what)
		return telemetry.LogRelation(a, action, a.Name(), toAtom)
	}
}

func Call(a Actor, req any, what, to string) (any, error) {
	toAtom := gen.Atom(to)

	if reply, err := a.Call(toAtom, req); err != nil {
		return reply, logError(err, a.Name(), to, what)
	} else {
		action := fmt.Sprintf("request %s", what)
		return reply, telemetry.LogRelation(a, action, a.Name(), toAtom)
	}
}

func Trigger(a Actor, e gen.Atom) error {
	event := string(e)
	namespace := strings.Split(event, "_")[0]
	from := fmt.Sprintf("%s_producer", namespace)
	action := fmt.Sprintf("trigger %s", event)
	return telemetry.LogRelation(a, action, gen.Atom(from), a.Name())
}

func RegisterActor(a Actor, name string) {
	telemetry.RegisterActor(a, name)
}

func TelemetryScheme(a Actor) []byte {
	return telemetry.BuildScheme(a)
}

func logError(err error, from gen.Atom, to, what string) error {
	if err != nil {
		fmt.Printf("!!! Intercom unexpected error [%s-%s->%s]: %s\n", from, what, to, err)
	}
	return err
}
