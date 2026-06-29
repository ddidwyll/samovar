package mqtt

import (
	"samovar/common/models"
	"samovar/lib/i"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	natiu "github.com/soypat/natiu-mqtt"

	"errors"
	"fmt"
)

type buffer = map[string]string

type client struct {
	act.Actor
	config        *config
	listenerAlias gen.Alias
	packetId      uint16
	topicBuffer   buffer
}

func newClient() gen.ProcessBehavior { return &client{} }

func (c *client) Init(args ...any) error {
	inter.RegisterActor(c, "[mqtt.client]")

	if err := c.initConfig(args...); err != nil {
		return err
	}

	c.resetBuffer()

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

var bufferedTopics = []string{"flag_otb", "otbor"}

func (c *client) resetBuffer() {
	c.topicBuffer = make(buffer, len(bufferedTopics))
}

func (c *client) createListener() error {
	listener := newListener(c.config)
	opts := gen.MetaOptions{}
	if alias, err := c.SpawnMeta(listener, opts); err != nil {
		return err
	} else {
		c.listenerAlias = alias
	}
	return nil
}

func (c *client) HandleMessage(_ gen.PID, msg any) error {
	m, ok := msg.(models.MqttMessage)
	if !ok {
		return errors.New("Unexpected mqtt message")
	}
	c.Log().Debug("mqtt.client.HandleMessage.msg.Topic: %v", m.Topic)
	if i.N(m.Topic, bufferedTopics...) {
		c.topicBuffer[m.Topic] = m.Text
		if len(c.topicBuffer) != len(bufferedTopics) {
			return nil
		}
		_, err := inter.Call(c, c.topicBuffer, "changes", "device_raw_state")
		if err == nil {
			c.resetBuffer()
			return nil
		} else {
			return err
		}
	}
	return inter.Send(c, m, "message", "mqtt_producer")
}

func (c *client) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	c.Log().Debug("mqtt.client.HandleCall.req: %v", req)
	msg, ok := req.(models.MqttMessage)
	if !ok {
		return nil, errors.New("Unexpected mqtt client change request")
	}
	client, err := c.Call(c.listenerAlias, "mqtt_client_request")
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
	c.Log().Debug("mqtt.client.publish.msg: %v", msg)
	pid := c.packetId + 1
	c.packetId = pid
	topic := fmt.Sprintf("/samovar/%s", msg.Topic)
	vars := natiu.VariablesPublish{[]byte(topic), pid}
	if flags, err := natiu.NewPublishFlags(0, false, false); err == nil {
		return client.PublishPayload(flags, vars, []byte(msg.Text))
	} else {
		return err
	}
}
