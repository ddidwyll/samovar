package mqtt

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"ergo.services/ergo/gen"
	natiu "github.com/soypat/natiu-mqtt"
)

type Listener struct {
	gen.MetaProcess
	config *config
	client *natiu.Client
	vconn  *natiu.VariablesConnect
	vsub   *natiu.VariablesSubscribe
}

func newListener(cfg *config) gen.MetaBehavior {
	var vconn natiu.VariablesConnect
	var vsub natiu.VariablesSubscribe

	return &Listener{vconn: &vconn, vsub: &vsub, config: cfg}
}

func (l *Listener) Init(process gen.MetaProcess) error {
	l.MetaProcess = process
	l.Log().Info("mqtt.Listener started (%v)", l.config)
	return nil
}

func (l *Listener) Start() error {
	l.createClient()

	if err := l.connect(); err != nil {
		l.Log().Error("mqtt Listener connect error: %s", err)
		return err
	}

	if err := l.subscribe(); err != nil {
		l.Log().Error("mqtt Listener subscribe error: %s", err)
		return err
	}

	for {
		err := withTimeout(func(ctx context.Context) error {
			return l.client.Ping(ctx)
		})
		if err == nil && !l.client.IsConnected() {
			err = errors.New("mqtt client not connected")
		}
		if err != nil {
			l.Log().Error("mqtt.Listener error: %s", err)
			return err
		}
		time.Sleep(time.Second)
	}
}

func (l *Listener) createClient() {
	onPub := func(_ natiu.Header, vpub natiu.VariablesPublish, r io.Reader) error {
		if message, err := io.ReadAll(r); err == nil {
			l.Log().Info("mqtt topic: %s", vpub.TopicName)
			l.Send(l.ID(), message)
			return nil
		} else {
			return err
		}
	}

	l.client = natiu.NewClient(natiu.ClientConfig{OnPub: onPub})
}

func (l *Listener) connect() error {
	const url = "127.0.0.1:1883"

	conn, err := net.Dial("tcp", url)

	if err != nil {
		return err
	}

	l.vconn.SetDefaultMQTT([]byte("samovar_client"))

	return withTimeout(func(ctx context.Context) error {
		return l.client.Connect(ctx, conn, l.vconn)
	})
}

func (l *Listener) subscribe() error {
	subscribeRequest := natiu.SubscribeRequest{[]byte("samovar"), natiu.QoS2}
	subscribeRequests := []natiu.SubscribeRequest{subscribeRequest}
	l.vsub.TopicFilters = subscribeRequests
	l.vsub.PacketIdentifier = 1

	return withTimeout(func(ctx context.Context) error {
		return l.client.Subscribe(ctx, *l.vsub)
	})
}

func (l *Listener) HandleCall(_ gen.PID, _ gen.Ref, request any) (any, error) {
	l.Log().Info("mqtt.Listener receive call request: %v", request)
	return nil, nil
}

func (l *Listener) HandleMessage(_ gen.PID, message any) error {
	l.Log().Info("mqtt.Listener receive message: %s", message)
	return nil
}

func (l *Listener) HandleInspect(_ gen.PID, items ...string) map[string]string {
	l.Log().Info("mqtt.Listener receive inspect request: %v", items)
	return nil
}

func (l *Listener) Terminate(reason error) {
	l.Log().Error("mqtt.Listener terminated (%s)", reason)
}

func withTimeout(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err := f(ctx)
	cancel()
	return err
}
