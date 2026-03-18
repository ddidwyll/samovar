package api

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"errors"
)

const sseStateProcessName = gen.Atom("sseState")

func factory_sseState() gen.ProcessBehavior {
	return &sseState{
		connections: make(map[gen.Alias]bool),
	}
}

type sseState struct {
	act.Actor
	connections map[gen.Alias]bool
	counter     uint
}

type sseStateReqAddConn struct{ alias gen.Alias }
type sseStateReqDelConn struct{ alias gen.Alias }
type sseStateReqCountIncr struct{}
type sseStateReqGet struct{ key string }

func (ss *sseState) Init(_ ...any) (err error) {
	ss.Log().Info("sseState started (%s)", ss.Name())
	return
}

func (ss *sseState) HandleCall(_ gen.PID, _ gen.Ref, request any) (any, error) {
	req, ok := request.(sseStateReqGet)

	if !ok {
		return nil, errors.New("sseState: unexpected call request")
	}

	switch req.key {
	case "counter":
		return ss.counter, nil
	case "connections":
		return ss.connections, nil
	default:
		return nil, errors.New("sseState: unexpected state key")
	}
}

func (ss *sseState) HandleMessage(_ gen.PID, message any) error {
	switch m := message.(type) {
	case sseStateReqAddConn:
		ss.connections[m.alias] = true
	case sseStateReqDelConn:
		delete(ss.connections, m.alias)
	case sseStateReqCountIncr:
		ss.counter++
	default:
		return errors.New("sseState: unexpected message")
	}

	return nil
}
