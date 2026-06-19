package mqtt

import (
	"samovar/lib/inter"
	"samovar/common/models"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	natiu "github.com/soypat/natiu-mqtt"

	"errors"
	"fmt"
)

type client struct {
	act.Actor
	config   *config
	listener gen.Alias
}

func newClient() gen.ProcessBehavior { return &client{} }

func (c *client) Init(args ...any) error {
	inter.RegisterActor(c, "[mqtt.client]")

	if err := c.initConfig(args...); err != nil {
		return err
	}

	return c.createListener()
}

func (c *client) initConfig(args ...any) error {
	if cfg, ok := args[0].(*config); !ok {
		return errors.New("invalid mqtt client config")
	} else {
		c.config = cfg
		return nil
	}
}

func (c *client) createListener() error {
	listener := newListener(c.config)
	opts := gen.MetaOptions{}
	if alias, err := c.SpawnMeta(listener, opts); err != nil {
		return err
	} else {
  	c.listener = alias
	}
	return nil
}

func (c *client) HandleMessage(_ gen.PID, msg any) error {
	c.Log().Debug("mqtt.client.HandleMessage.msg: %v", msg)
	return inter.Send(c, msg, "message", "mqtt_producer")
}

func (c *client) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	c.Log().Info("mqtt.client.HandleCall.req: %v", req)
	msg, ok := req.(models.MqttMessage)
	if !ok {
  	return nil, errors.New("Unexpected mqtt client change request")
	}
	client, err := c.Call(c.listener, "mqtt_client_request")
	if err != nil {
  	return nil, err
	}
	mqttClient, ok := client.(*natiu.Client)
	if !ok {
  	return nil, errors.New("Unexpected mqtt listener client")
	}
	return "ok", c.publish(mqttClient, msg)
}

func (c *client) publish(client *natiu.Client, msg models.MqttMessage) error {
  topic := fmt.Sprintf("/samovar/%s", msg.Topic)
  vars := natiu.VariablesPublish{[]byte(topic), 1}
  if flags, err := natiu.NewPublishFlags(0, false, false); err == nil {
    return client.PublishPayload(flags, vars, []byte(msg.Text))
  } else {
    return err
  }
}
