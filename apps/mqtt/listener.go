package mqtt

import (
	"context"
	"io"
	"net"
	"time"
	"errors"
	// "ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	natiu "github.com/soypat/natiu-mqtt"
)

type Listener struct {
	gen.MetaProcess
	client *natiu.Client
	target gen.PID
}

func newListener() gen.MetaBehavior {
	return &Listener{}
}

func (l *Listener) Init(process gen.MetaProcess) error {
	l.MetaProcess = process
	l.target = process.Parent()

	var config natiu.ClientConfig

	// buf := make([]byte, 1500)
	// config.Decoder = natiu.DecoderNoAlloc{buf}

	config.OnPub = func(_ natiu.Header, _ natiu.VariablesPublish, r io.Reader) error {
		message, _ := io.ReadAll(r)
		l.Log().Info("mqtt message: %s", message)
		l.Send(l.target, message)
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
		return err
	}

	var varConn natiu.VariablesConnect
	varConn.SetDefaultMQTT([]byte("samovar_client"))
	err = withTimeout(5, func(ctx context.Context) error {
  	return l.client.Connect(ctx, conn, &varConn)
	})

	if err != nil {
		l.Log().Error("mqtt Listener error: %s", err)
		return err
	}

	subscribeRequest := natiu.SubscribeRequest{[]byte("topic"), natiu.QoS2}
	subscribeRequests := []natiu.SubscribeRequest{subscribeRequest}
	varSub := natiu.VariablesSubscribe{subscribeRequests, 1}
	err = withTimeout(5, func(ctx context.Context) error {
  	return l.client.Subscribe(ctx, varSub)
	})

	if err != nil {
		l.Log().Error("mqtt Listener error: %s", err)
		return err
	}

	i := 0
	for {
  	l.Log().Info("mqtt loop!%d", i)
		err := withTimeout(2, func(ctx context.Context) error {
      return l.client.Ping(ctx)
		})
  	if err == nil && !l.client.IsConnected() {
  	  err = errors.New("mqtt client not connected")
  	}
		if err != nil {
			l.Log().Error("mqtt.Listener error: %s", err)
			return err
		}
		time.Sleep(5*time.Second)
		i++
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

func withTimeout(seconds time.Duration, f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), seconds*time.Second)
	err := f(ctx)
	cancel()
	return err
}

