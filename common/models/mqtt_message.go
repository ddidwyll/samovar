package models

import (
	"path"
	"time"
)

type MqttMessage struct {
	Topic     string
	Text      string
	Timestamp int64
}

func NewMqttMessage(topicPath, text []byte) MqttMessage {
	_, topic := path.Split(string(topicPath))
	timestamp := time.Now().UnixMilli()

	return MqttMessage{topic, string(text), timestamp}
}
