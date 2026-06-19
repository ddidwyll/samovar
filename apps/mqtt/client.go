package mqtt

import (
	"samovar/common/models"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	natiu "github.com/soypat/natiu-mqtt"

	"context"
	"errors"
	"io"
	"net"
	"time"
)

type client struct {
	act.Actor
	config *config
	mqtt   *natiu.Client
	// listener gen.Alias
}

func newClient() gen.ProcessBehavior { return &client{} }

func (c *client) Init(args ...any) error {
	inter.RegisterActor(c, "[mqtt.client]")

	if err := c.initConfig(args...); err != nil {
		return err
	}

	return c.start()
}

func (c *client) initConfig(args ...any) error {
	if cfg, ok := args[0].(*config); !ok {
		return errors.New("invalid mqtt client config")
	} else {
		c.config = cfg
		return nil
	}
}

func (c *client) start() (err error) {
	c.createClient()
	if err = c.connect(); err != nil {
		return err
	}
	if err = c.subscribe(); err != nil {
		return err
	}
	if err = c.createListener(); err != nil {
		return err
	}
	return nil
}

func (c *client) createClient() {
	onPub := func(_ natiu.Header, vpub natiu.VariablesPublish, r io.Reader) error {
		if text, err := io.ReadAll(r); err == nil {
			msg := models.NewMqttMessage(vpub.TopicName, text)
			c.Log().Info("mqtt.client: %v, %v", c.State(), msg)
			// return inter.Send(c, msg, "message", "mqtt_producer")
			return nil
		} else {
			return err
		}
	}

	cfg := natiu.ClientConfig{OnPub: onPub}
	c.mqtt = natiu.NewClient(cfg)
}

func (c *client) connect() error {
	var vconn natiu.VariablesConnect

	url := c.config.url()
	conn, err := net.Dial("tcp", url)

	if err != nil {
		return err
	}

	clientID := c.config.clientID()
	vconn.SetDefaultMQTT(clientID)

	return c.withTimeout(func(ctx context.Context) error {
		return c.mqtt.Connect(ctx, conn, &vconn)
	})
}

func (c *client) subscribe() error {
	var vsub natiu.VariablesSubscribe

	topic := c.config.topic()
	topicFilter := natiu.SubscribeRequest{topic, natiu.QoS2}
	vsub.TopicFilters = []natiu.SubscribeRequest{topicFilter}
	vsub.PacketIdentifier = 1

	return c.withTimeout(func(ctx context.Context) error {
		return c.mqtt.Subscribe(ctx, vsub)
	})
}

func (c *client) createListener() error {
	listener := newListener(c.mqtt)
	_, err := c.SpawnMeta(listener, gen.MetaOptions{})
	if err != nil {
		return err
	}
	// c.listener = alias
	return nil
}

func (c *client) withTimeout(f func(context.Context) error) error {
	timeout := c.config.Timeout * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	err := f(ctx)
	cancel()
	return err
}

func (c *client) HandleMessage(_ gen.PID, msg any) error {
	// c.Log().Info("mqtt.client receive message: %#v", msg)
	return nil
	// return inter.Send(c, msg, "message", "mqtt_producer")
}
