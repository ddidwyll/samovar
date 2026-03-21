package api

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type sup struct{ act.Supervisor }

func newSup() gen.ProcessBehavior { return &sup{} }

func (s *sup) Init(args ...any) (spec act.SupervisorSpec, err error) {
	spec.Type = act.SupervisorTypeOneForOne
	spec.Restart.Strategy = act.SupervisorStrategyTransient
	spec.Restart.Intensity = 2
	spec.Restart.Period = 5

	spec.Children = []act.SupervisorChildSpec{
		{
			Name:    "api_sse_state",
			Factory: newsseState,
		},
		{
			Name:    "api_web",
			Factory: createWeb,
		},
		{
			Name:    "api_actor",
			Factory: createActor,
		},
	}

	return spec, err
}

// HandleChildStart invoked on a successful child process starting if option EnableHandleChild
// was enabled in act.SupervisorSpec
func (s *sup) HandleChildStart(name gen.Atom, pid gen.PID) error {
	return nil
}

// HandleChildTerminate invoked on a child process termination if option EnableHandleChild
// was enabled in act.SupervisorSpec
func (s *sup) HandleChildTerminate(name gen.Atom, pid gen.PID, reason error) error {
	return nil
}

// HandleMessage invoked if Supervisor received a message sent with gen.Process.Send(...).
// Non-nil value of the returning error will cause termination of this process.
// To stop this process normally, return gen.TerminateReasonNormal or
// gen.TerminateReasonShutdown. Any other - for abnormal termination.
func (s *sup) HandleMessage(from gen.PID, message any) error {
	s.Log().Info("supervisor got message from %s", from)
	return nil
}

// HandleCall invoked if Supervisor got a synchronous request made with gen.Process.Call(...).
// Return nil as a result to handle this request asynchronously and
// to provide the result later using the gen.Process.SendResponse(...) method.
func (s *sup) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	s.Log().Info("supervisor got request from %s with reference %s", from, ref)
	return gen.Atom("pong"), nil
}

// Terminate invoked on a termination supervisor process
func (s *sup) Terminate(reason error) {
	s.Log().Info("supervisor terminated with reason: %s", reason)
}

// HandleInspect invoked on the request made with gen.Process.Inspect(...)
func (s *sup) HandleInspect(from gen.PID, item ...string) map[string]string {
	s.Log().Info("supervisor got inspect request from %s", from)
	return nil
}
