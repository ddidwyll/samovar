package client

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func newState() gen.ProcessBehavior { return &state{} }

type state struct{ act.Actor }

// Init invoked on s start this process.
func (s *state) Init(args ...any) error {
	s.Log().Debug("started process with name %s and args %v", s.Name(), args)
	return nil
}

//
// Methods below are optional, so you can remove those that aren't be used
//

// HandleMessage invoked if state received s message sent with gen.Process.Send(...).
// Non-nil value of the returning error will cause termination of this process.
// To stop this process normally, return gen.TerminateReasonNormal
// or any other for abnormal termination.
func (s *state) HandleMessage(from gen.PID, message any) error {
	s.Log().Debug("got message from %s", from)
	return nil
}

// HandleCall invoked if state got s synchronous request made with gen.Process.Call(...).
// Return nil as s result to handle this request asynchronously and
// to provide the result later using the gen.Process.SendResponse(...) method.
func (s *state) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	s.Log().Debug("got request from %s with reference %s", from, ref)
	return gen.Atom("pong"), nil
}

// Terminate invoked on s termination process
func (s *state) Terminate(reason error) {
	s.Log().Debug("terminated with reason: %s", reason)
}

// HandleMessageName invoked if split handling was enabled using SetSplitHandle(true)
// and message has been sent by name
func (s *state) HandleMessageName(name gen.Atom, from gen.PID, message any) error {
	return nil
}

// HandleMessageAlias invoked if split handling was enabled using SetSplitHandle(true)
// and message has been sent by alias
func (s *state) HandleMessageAlias(alias gen.Alias, from gen.PID, message any) error {
	return nil
}

// HandleCallName invoked if split handling was enabled using SetSplitHandle(true)
// and request was made by name
func (s *state) HandleCallName(name gen.Atom, from gen.PID, ref gen.Ref, request any) (any, error) {
	return gen.Atom("pong"), nil
}

// HandleCallAlias invoked if split handling was enabled using SetSplitHandle(true)
// and request was made by alias
func (s *state) HandleCallAlias(alias gen.Alias, from gen.PID, ref gen.Ref, request any) (any, error) {
	return gen.Atom("pong"), nil
}

// HandleLog invoked on s log message if this process was added as s logger.
// See https://docs.ergo.services/basics/logging for more information
func (s *state) HandleLog(message gen.MessageLog) error {
	return nil
}

// HandleEvent invoked on an event message if this process got subscribed on
// this event using gen.Process.LinkEvent or gen.Process.MonitorEvent
// See https://docs.ergo.services/basics/events for more information
func (s *state) HandleEvent(message gen.MessageEvent) error {
	return nil
}

// HandleInspect invoked on the request made with gen.Process.Inspect(...)
func (s *state) HandleInspect(from gen.PID, item ...string) map[string]string {
	s.Log().Debug("got inspect request from %s", from)
	return nil
}
