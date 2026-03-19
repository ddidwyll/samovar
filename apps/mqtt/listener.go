package mqtt

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Listener struct {
	act.Actor
}

func factory_Listener() gen.ProcessBehavior {
	return &Listener{}
}

func (a *Listener) Init(_ ...any) error {
	a.Log().Info("mqtt.Listener started (%s)", a.Name())
	return nil
}
