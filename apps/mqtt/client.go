package mqtt

import (
	"samovar/common/models"
	"samovar/lib/inter"
	"samovar/lib/i"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	natiu "github.com/soypat/natiu-mqtt"

	"errors"
	"fmt"
)

type buffer = map[string]string

type client struct {
	act.Actor
	config   *config
	listenerAlias gen.Alias
	packetId uint16
	inputBuffer buffer
}

var topics = []string{
  "term_d",
  "term_c",
  "term_k",
  "power",
  "term_v",
  "press_a",
  "term_c_max",
  "term_c_min",
  "term_k_m",
  "term_d_m",
  "term_nasos",
  "power_m",
  "term_vent",
  "flag_otb",
  "otbor",
  "min_otb",
  "sek_otb",
  "time_stop",
  "otbor_minus",
  "otbor_g_1",
  "otbor_g_2",
  "otbor_t",
  "delta_t",
  "count_vent",
  "num_error",
  "kontaktor",
  "last_ping",
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

func (c *client) resetBuffer() {
  c.inputBuffer = make(buffer, len(topics))
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
	c.Log().Info("mqtt.client.HandleMessage.msg.Topic: %v", m.Topic)
	if !ok {
  	return errors.New("Unexpected mqtt message")
	}
  if !i.N(m.Topic, topics...) {
    if i.N(m.Topic, "last_ping") {
    	return inter.Send(c, m, "message", "mqtt_producer")
    } else {
      return nil
    }
  }
  c.inputBuffer[m.Topic] = m.Text
  if len(c.inputBuffer) != len(topics) {
    return nil
  }
  _, err := inter.Call(c, c.inputBuffer, "changes", "device_raw_state")
	if err == nil {
  	c.Log().Info("mqtt.client.HandleMessage: send package")
    c.resetBuffer()
    return nil
	} else {
  	return err
	}
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
