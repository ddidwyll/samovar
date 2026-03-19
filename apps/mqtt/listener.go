package mqtt

import (
	"io"
	"net"
	"context"
	"time"
	// "ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	natiu "github.com/soypat/natiu-mqtt"
)

type Listener struct {
	gen.MetaProcess
	client *natiu.Client
}

func newListener() gen.MetaBehavior {
	return &Listener{}
}

func (l *Listener) Init(process gen.MetaProcess) error {
  l.MetaProcess = process

	var config natiu.ClientConfig

	buf := make([]byte, 1500)

	config.Decoder = natiu.DecoderNoAlloc{buf}

	config.OnPub = func(_ natiu.Header, _ natiu.VariablesPublish, r io.Reader) error {
		message, _ := io.ReadAll(r)
		l.Log().Info("mqtt message: %s", message)
		l.Send(l.Parent(), message)
		return nil
	}

	l.client = natiu.NewClient(config)

	l.Log().Info("mqtt.Listener started (%s)", l.ID())
	return nil
}

func (l *Listener) Start() error {
  const url = "127.0.0.1:1883"

  conn, err := net.Dial("tcp", url)

  if err != nil {
    l.Log().Error("mqtt Listener error: %s", err)
    panic(err)
  }

  var varConn natiu.VariablesConnect
  varConn.SetDefaultMQTT([]byte("salamanca"))
  ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
  err = l.client.Connect(ctx, conn, &varConn)
  cancel()

  if err != nil {
    l.Log().Error("mqtt Listener error: %s", err)
    panic(err)
  }

  for {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    pingErr := l.client.Ping(ctx)
    cancel()
    if pingErr != nil {
      l.Log().Error("mqtt.Listener error: %s", pingErr)
      l.Log().Error("mqtt.Listener error: %s", l.client.Err())
      panic(pingErr)
    }
    l.Log().Info("mqtt ping!")
  }
}

func (l *Listener) HandleCall(_ gen.PID, _ gen.Ref, request any) (any, error) {
  l.Log().Info("mqtt.Listener receive call request: %v", request)
  return nil, nil
}

func (l *Listener) HandleMessage(_ gen.PID, message any) error {
  l.Log().Info("mqtt.Listener receive message: %v", message)
  return nil
}

func (l *Listener) HandleInspect(_ gen.PID, items ...string) map[string]string {
  l.Log().Info("mqtt.Listener receive inspect request: %v", items)
  return nil
}

func (l *Listener) Terminate(reason error) {
  l.Log().Error("mqtt.Listener terminated (%s)", reason)
}
