package bus

import "samovar/apps/mqtt"

type MqttNewMessage struct {
	mqtt.Message
}
