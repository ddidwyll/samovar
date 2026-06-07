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
		return err
	} else {
		action := fmt.Sprintf("send %s", what)
		// 	fmt.Printf("[%s] >> %s >> [%s]\n", a.Name(), action, to)
		return telemetry.LogRelation(a, action, a.Name(), toAtom)
	}
}

func Call(a Actor, req any, what, to string) (any, error) {
	toAtom := gen.Atom(to)

	if reply, err := a.Call(toAtom, req); err != nil {
		return reply, err
	} else {
		action := fmt.Sprintf("request %s", what)
		// 	fmt.Printf("[%s] >> %s >> [%s]\n", a.Name(), action, to)
		return reply, telemetry.LogRelation(a, action, a.Name(), toAtom)
	}
}

func Trigger(a Actor, event gen.Atom, from string) error {
	fromAtom := gen.Atom(from)
	action := fmt.Sprintf("trigger %s", event)
	// fmt.Printf("[%s] >> %s >> [%s]\n", from, action, a.Name())
	return telemetry.LogRelation(a, action, fromAtom, a.Name())
}

func RegisterActor(a Actor, name string) {
	// fmt.Printf("<< %s | %s >>\n", a.Name(), name)
	telemetry.RegisterActor(a, name)
}
