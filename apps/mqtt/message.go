package mqtt

import (
	"time"
	"path"
)

type Message struct {
	Topic     string
	Text      string
	Timestamp int64
}

func newMessage(topicPath, text []byte) Message {
  _, topic := path.Split(string(topicPath))
  timestamp := time.Now().UnixMicro()

	return Message{topic, string(text), timestamp}
}
