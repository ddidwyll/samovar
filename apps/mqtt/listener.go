package mqtt

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"samovar/common/models"

	"ergo.services/ergo/gen"
	natiu "github.com/soypat/natiu-mqtt"
)

type listener struct {
	gen.MetaProcess
	config *config
	client *natiu.Client
	vconn  *natiu.VariablesConnect
	vsub   *natiu.VariablesSubscribe
}

func newListener(cfg *config) gen.MetaBehavior {
	var vconn natiu.VariablesConnect
	var vsub natiu.VariablesSubscribe

	return &listener{vconn: &vconn, vsub: &vsub, config: cfg}
}

func (l *listener) Init(process gen.MetaProcess) error {
	l.MetaProcess = process
	l.createClient()
	l.Log().Debug("mqtt.listener started (%v)", l.config)
	return nil
}

func (l *listener) Start() error {
	if err := l.connect(); err != nil {
		l.Log().Error("mqtt.listener connect error: %s", err)
		return err
	}

	if err := l.subscribe(); err != nil {
		l.Log().Error("mqtt.listener subscribe error: %s", err)
		return err
	}

	for {
		err := withTimeout(func(ctx context.Context) error {
			return l.client.Ping(ctx)
		})
		if err == nil && !l.client.IsConnected() {
			err = errors.New("mqtt.listener error: not connected")
		}
		if err != nil {
			l.Log().Error("mqtt.listener error: %s", err)
			return err
		}
		time.Sleep(time.Second)
	}
}

func (l *listener) createClient() {
	onPub := func(_ natiu.Header, vpub natiu.VariablesPublish, r io.Reader) error {
		if text, err := io.ReadAll(r); err == nil {
			msg := models.NewMqttMessage(vpub.TopicName, text)
			l.Send(l.Parent(), msg)
			return nil
		} else {
			return err
		}
	}

	l.client = natiu.NewClient(natiu.ClientConfig{OnPub: onPub})
}

func (l *listener) connect() error {
	url := l.config.url()
	conn, err := net.Dial("tcp", url)

	if err != nil {
		return err
	}

	clientID := l.config.clientID()
	l.vconn.SetDefaultMQTT(clientID)

	return withTimeout(func(ctx context.Context) error {
		return l.client.Connect(ctx, conn, l.vconn)
	})
}

func (l *listener) subscribe() error {
	topic := l.config.topic()
	topicFilter := natiu.SubscribeRequest{topic, natiu.QoS2}
	l.vsub.TopicFilters = []natiu.SubscribeRequest{topicFilter}
	l.vsub.PacketIdentifier = 1

	return withTimeout(func(ctx context.Context) error {
		return l.client.Subscribe(ctx, *l.vsub)
	})
}

func (l *listener) HandleCall(_ gen.PID, _ gen.Ref, request any) (any, error) {
	l.Log().Debug("mqtt.listener receive call request: %v", request)
	return nil, nil
}

func (l *listener) HandleMessage(_ gen.PID, msg any) error {
	l.Log().Debug("mqtt.listener receive message: %s", msg)
	return nil
}

func (l *listener) HandleInspect(_ gen.PID, items ...string) map[string]string {
	l.Log().Debug("mqtt.listener receive inspect request: %v", items)
	return nil
}

func (l *listener) Terminate(reason error) {
	l.Log().Debug("mqtt.listener terminated (%s)", reason)
}

func withTimeout(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := f(ctx)
	cancel()
	return err
}
