package inter

import (
	"samovar/apps/telemetry"

	"ergo.services/ergo/gen"

	"fmt"
)

type Actor = telemetry.Actor

func Send(a Actor, msg any, what, to string) error {
	toAtom := gen.Atom(to)

	if err := a.Send(toAtom, msg); err != nil {
		return logError(err)
	} else {
		action := fmt.Sprintf("send %s", what)
		return logError(telemetry.LogRelation(a, action, a.Name(), toAtom))
	}
}

func Call(a Actor, req any, what, to string) (any, error) {
	toAtom := gen.Atom(to)

	if reply, err := a.Call(toAtom, req); err != nil {
		return reply, logError(err)
	} else {
		action := fmt.Sprintf("request %s", what)
		return reply, logError(telemetry.LogRelation(a, action, a.Name(), toAtom))
	}
}

func Trigger(a Actor, event, producer gen.Atom) error {
	action := fmt.Sprintf("trigger %s", string(event))
	return logError(telemetry.LogRelation(a, action, producer, a.Name()))
}

func RegisterActor(a Actor, name string) {
	telemetry.RegisterActor(a, name)
}

func TelemetryScheme(a Actor) []byte {
	return telemetry.BuildScheme(a)
}

func logError(err error) error {
	if err != nil {
		fmt.Printf("!!! Intercom unexpected error: %s\n", err)
	}
	return err
}
