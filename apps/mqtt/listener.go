package mqtt

import (
	"samovar/common/models"

	"ergo.services/ergo/gen"

	natiu "github.com/soypat/natiu-mqtt"

	"context"
	"io"
	"net"
	"time"
)

type listener struct {
	gen.MetaProcess
	config *config
	client *natiu.Client
}

func newListener(cfg *config) gen.MetaBehavior {
	return &listener{config: cfg}
}

func (l *listener) Init(process gen.MetaProcess) error {
	l.MetaProcess = process
	return nil
}

func (l *listener) Start() (err error) {
	l.createClient()
	if err = l.connect(); err != nil {
		return err
	}
	if err = l.subscribe(); err != nil {
		return err
	}
	restart := func(reason error) {
		l.Log().Error("mqtt.listener.error: %s", err)
		l.Start()
	}
	for {
		err = l.withTimeout(func(ctx context.Context) error {
			return l.client.Ping(ctx)
		})
		// err = l.client.HandleNext()
		if err != nil {
			break
		} else {
			if err = l.sendTimestamp(); err != nil {
				break
			} else {
				time.Sleep(250 * time.Millisecond)
			}
		}
	}
	restart(err)
	return nil
}

func (l *listener) sendTimestamp() error {
	sendTs := func(k string, t time.Time) error {
		ts := t.Format(time.TimeOnly)
		msg := models.NewMqttMessage([]byte(k), []byte(ts))
		return l.Send(l.Parent(), msg)
	}
	// if err := sendTs("last_tx", l.client.LastTx()); err != nil {
	//  	return err
	// }
	// if err := sendTs("last_rx", l.client.LastRx()); err != nil {
	//  	return err
	// }
	if err := sendTs("last_ping", time.Now()); err != nil {
		return err
	}
	return nil
}
func (l *listener) onPub(_ natiu.Header, vpub natiu.VariablesPublish, r io.Reader) error {
	if text, err := io.ReadAll(r); err == nil {
		msg := models.NewMqttMessage(vpub.TopicName, text)
		l.Log().Debug("mqtt.listener.onPub.msg: %v", msg)
		return l.Send(l.Parent(), msg)
	} else {
		return err
	}
}

func (l *listener) createClient() {
	cfg := natiu.ClientConfig{OnPub: l.onPub}
	l.client = natiu.NewClient(cfg)
}

func (l *listener) connect() error {
	var vconn natiu.VariablesConnect

	url := l.config.url()
	conn, err := net.Dial("tcp", url)

	if err != nil {
		return err
	}

	clientID := l.config.clientID()
	vconn.SetDefaultMQTT(clientID)

	return l.withTimeout(func(ctx context.Context) error {
		return l.client.Connect(ctx, conn, &vconn)
	})
}

func (l *listener) subscribe() error {
	var vsub natiu.VariablesSubscribe

	topic := l.config.topic()
	topicFilter := natiu.SubscribeRequest{topic, natiu.QoS2}
	vsub.TopicFilters = []natiu.SubscribeRequest{topicFilter}
	vsub.PacketIdentifier = 1

	return l.withTimeout(func(ctx context.Context) error {
		return l.client.Subscribe(ctx, vsub)
	})
}

func (l *listener) withTimeout(fn func(context.Context) error) error {
	timeout := l.config.Timeout * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()
	return fn(ctx)
}

func (l *listener) Terminate(reason error) {
	l.Log().Debug("mqtt.listener.Terminate.reason: %s", reason)
}

func (l *listener) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	l.Log().Debug("mqtt.listener.HandleCall.req: %v", req)
	return l.client, nil

}
func (l *listener) HandleMessage(_ gen.PID, msg any) error {
	l.Log().Debug("mqtt.listener.HandleMessage.msg: %s", msg)
	return nil
}

func (l *listener) HandleInspect(_ gen.PID, items ...string) map[string]string {
	l.Log().Debug("mqtt.listener.HandleInspect.items: %v", items)
	return nil
}
