package mqtt

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_MqttListener() gen.ProcessBehavior {
  return &MqttListener{}
}

type MqttListener struct {
  act.Actor
  port byte
  host string
}
