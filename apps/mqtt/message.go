package mqtt

import (
	"time"
)

type message struct {
	topic     string
	text      string
	timestamp int64
}

func newMessage(topic, text []byte) message {
	return message{
		topic:     string(topic),
		text:      string(text),
		timestamp: time.Now().UnixMicro(),
	}
}
